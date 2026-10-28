package acp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

// The end-to-end tests run a fake agent loop as a real subprocess (this test
// binary, re-executed) and drive it over stdio with a fake messaging client:
// the shape of a chat bridge talking to `some-agent acp`.

const fakeAgentEnv = "ACP_FAKE_AGENT"

func TestMain(m *testing.M) {
	if os.Getenv(fakeAgentEnv) == "1" {
		agent := &loopAgent{conn: NewAgentConn(os.Stdin, os.Stdout)}
		agent.conn.Serve(context.Background(), agent)
		<-agent.conn.Done()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// loopAgent is a fake agent loop: per turn it thinks, plans, runs one tool
// behind a permission request, and streams its answer.
type loopAgent struct {
	conn *AgentConn
}

func (a *loopAgent) Initialize(_ context.Context, req InitializeRequest) (InitializeResponse, error) {
	return InitializeResponse{ProtocolVersion: req.ProtocolVersion, AgentCapabilities: &AgentCapabilities{LoadSession: true}}, nil
}

func (a *loopAgent) Authenticate(context.Context, AuthenticateRequest) (AuthenticateResponse, error) {
	return AuthenticateResponse{}, nil
}

func (a *loopAgent) NewSession(context.Context, NewSessionRequest) (NewSessionResponse, error) {
	return NewSessionResponse{SessionID: "chat-1"}, nil
}

func (a *loopAgent) Cancel(context.Context, CancelNotification) error { return nil }

// LoadSession replays the stored conversation before answering, per spec.
func (a *loopAgent) LoadSession(ctx context.Context, req LoadSessionRequest) (LoadSessionResponse, error) {
	if err := a.send(ctx, req.SessionID, SessionUpdate{UserMessageChunk: text("earlier question")}); err != nil {
		return LoadSessionResponse{}, err
	}
	return LoadSessionResponse{}, a.send(ctx, req.SessionID, SessionUpdate{AgentMessageChunk: text("earlier answer")})
}

func (a *loopAgent) SetSessionMode(ctx context.Context, req SetSessionModeRequest) (SetSessionModeResponse, error) {
	return SetSessionModeResponse{}, a.send(ctx, req.SessionID, SessionUpdate{CurrentModeUpdate: &CurrentModeUpdate{CurrentModeID: req.ModeID}})
}

func (a *loopAgent) Prompt(ctx context.Context, req PromptRequest) (PromptResponse, error) {
	id := req.SessionID
	steps := []SessionUpdate{
		{AgentThoughtChunk: text("I should list the files.")},
		{Plan: &Plan{Entries: []PlanEntry{{Content: "list files", Priority: PlanEntryPriorityHigh, Status: PlanEntryStatusInProgress}}}},
		{ToolCall: &ToolCall{ToolCallID: "t1", Title: "ls", Kind: ToolKindExecute, Status: ToolCallStatusPending}},
	}
	for _, u := range steps {
		if err := a.send(ctx, id, u); err != nil {
			return PromptResponse{}, err
		}
	}
	perm, err := a.conn.RequestPermission(ctx, RequestPermissionRequest{SessionID: id, ToolCall: ToolCallUpdate{ToolCallID: "t1"},
		Options: []PermissionOption{{OptionID: "allow", Name: "Allow", Kind: PermissionOptionKindAllowOnce}, {OptionID: "deny", Name: "Deny", Kind: PermissionOptionKindRejectOnce}}})
	if err != nil {
		return PromptResponse{}, err
	}
	if perm.Outcome.Cancelled != nil {
		return PromptResponse{StopReason: StopReasonCancelled}, nil
	}
	status, answer := ToolCallStatusFailed, "I was not allowed to list the files."
	if perm.Outcome.Selected.OptionID == "allow" {
		status, answer = ToolCallStatusCompleted, "There are two files: a.go and b.go."
	}
	if err := a.send(ctx, id, SessionUpdate{ToolCallUpdate: &ToolCallUpdate{ToolCallID: "t1", Status: status}}); err != nil {
		return PromptResponse{}, err
	}
	for _, word := range strings.SplitAfter(answer, " ") {
		if err := a.send(ctx, id, SessionUpdate{AgentMessageChunk: text(word)}); err != nil {
			return PromptResponse{}, err
		}
	}
	return PromptResponse{StopReason: StopReasonEndTurn}, nil
}

func (a *loopAgent) send(ctx context.Context, id SessionID, u SessionUpdate) error {
	return a.conn.SessionUpdate(ctx, SessionNotification{SessionID: id, Update: u})
}

func text(s string) *ContentChunk {
	return &ContentChunk{Content: ContentBlock{Text: &TextContent{Text: s}}}
}

// chatClient is a fake messaging client: it renders updates as chat lines the
// way a Telegram or Signal bridge would, and answers permission prompts as a
// user tapping a button. A nil choice never answers, like a user who walked off.
type chatClient struct {
	mu     sync.Mutex
	lines  []string
	reply  strings.Builder
	choice *PermissionOptionID
	asked  chan struct{}
}

func (c *chatClient) SessionUpdate(_ context.Context, n SessionNotification) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	u := n.Update
	if u.AgentMessageChunk != nil {
		c.reply.WriteString(u.AgentMessageChunk.Content.Text.Text)
		return nil
	}
	c.flush() // anything else ends the streamed reply
	switch {
	case u.UserMessageChunk != nil:
		c.lines = append(c.lines, "you: "+u.UserMessageChunk.Content.Text.Text)
	case u.AgentThoughtChunk != nil:
		c.lines = append(c.lines, "thinking: "+u.AgentThoughtChunk.Content.Text.Text)
	case u.Plan != nil:
		c.lines = append(c.lines, "plan: "+u.Plan.Entries[0].Content)
	case u.ToolCall != nil:
		c.lines = append(c.lines, "tool: "+u.ToolCall.Title)
	case u.ToolCallUpdate != nil:
		c.lines = append(c.lines, "tool "+string(u.ToolCallUpdate.Status))
	case u.CurrentModeUpdate != nil:
		c.lines = append(c.lines, "mode: "+string(u.CurrentModeUpdate.CurrentModeID))
	}
	return nil
}

// flush posts the streamed reply as one chat message, the way a bridge batches
// chunks instead of sending a message per word.
func (c *chatClient) flush() {
	if c.reply.Len() > 0 {
		c.lines = append(c.lines, "agent: "+c.reply.String())
		c.reply.Reset()
	}
}

func (c *chatClient) RequestPermission(ctx context.Context, req RequestPermissionRequest) (RequestPermissionResponse, error) {
	c.mu.Lock()
	c.lines = append(c.lines, fmt.Sprintf("approve? [%s] [%s]", req.Options[0].Name, req.Options[1].Name))
	c.mu.Unlock()
	if c.asked != nil {
		close(c.asked)
	}
	if c.choice == nil {
		<-ctx.Done()
		return RequestPermissionResponse{}, ctx.Err()
	}
	return RequestPermissionResponse{Outcome: RequestPermissionOutcome{Selected: &SelectedPermissionOutcome{OptionID: *c.choice}}}, nil
}

func (c *chatClient) transcript() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flush()
	return strings.Join(c.lines, "\n")
}

func spawnLoop(t *testing.T, client *chatClient) (*ClientConn, SessionID) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	cmd.Env = append(os.Environ(), fakeAgentEnv+"=1")
	cmd.Stderr = os.Stderr
	conn, err := Spawn(ctx, client, cmd)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if _, err := conn.Initialize(ctx, InitializeRequest{ProtocolVersion: V1}); err != nil {
		t.Fatal(err)
	}
	sess, err := conn.NewSession(ctx, NewSessionRequest{Cwd: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	return conn, sess.SessionID
}

func prompt(t *testing.T, conn *ClientConn, id SessionID) StopReason {
	t.Helper()
	resp, err := conn.Prompt(context.Background(), PromptRequest{SessionID: id, Prompt: []ContentBlock{{Text: &TextContent{Text: "list the files"}}}})
	if err != nil {
		t.Fatal(err)
	}
	return resp.StopReason
}

func TestE2EApprovedTurn(t *testing.T) {
	allow := PermissionOptionID("allow")
	client := &chatClient{choice: &allow}
	conn, id := spawnLoop(t, client)
	if stop := prompt(t, conn, id); stop != StopReasonEndTurn {
		t.Fatalf("stop reason = %q", stop)
	}
	want := "thinking: I should list the files.\nplan: list files\ntool: ls\napprove? [Allow] [Deny]\ntool completed\nagent: There are two files: a.go and b.go."
	if got := client.transcript(); got != want {
		t.Fatalf("chat:\n%s\nwant:\n%s", got, want)
	}
}

func TestE2EDeniedTurn(t *testing.T) {
	deny := PermissionOptionID("deny")
	client := &chatClient{choice: &deny}
	conn, id := spawnLoop(t, client)
	if stop := prompt(t, conn, id); stop != StopReasonEndTurn {
		t.Fatalf("stop reason = %q", stop)
	}
	if got := client.transcript(); !strings.Contains(got, "tool failed\nagent: I was not allowed to list the files.") {
		t.Fatalf("chat:\n%s", got)
	}
}

// The user cancels while the approval prompt is still open: the prompt
// resolves as cancelled on the client, and the turn ends cancelled.
func TestE2ECancelWhileAwaitingApproval(t *testing.T) {
	client := &chatClient{asked: make(chan struct{})}
	conn, id := spawnLoop(t, client)
	stop := make(chan StopReason, 1)
	go func() { stop <- prompt(t, conn, id) }()
	<-client.asked
	if err := conn.Cancel(context.Background(), CancelNotification{SessionID: id}); err != nil {
		t.Fatal(err)
	}
	if got := <-stop; got != StopReasonCancelled {
		t.Fatalf("stop reason = %q, want cancelled", got)
	}
}

func TestE2ELoadReplaysAndModeSwitches(t *testing.T) {
	client := &chatClient{}
	conn, id := spawnLoop(t, client)
	ctx := context.Background()
	if _, err := conn.LoadSession(ctx, LoadSessionRequest{SessionID: id, Cwd: t.TempDir()}); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.SetSessionMode(ctx, SetSessionModeRequest{SessionID: id, ModeID: "auto"}); err != nil {
		t.Fatal(err)
	}
	want := "you: earlier question\nagent: earlier answer\nmode: auto"
	if got := client.transcript(); got != want {
		t.Fatalf("chat:\n%s\nwant:\n%s", got, want)
	}
}
