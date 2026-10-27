package acp

import (
	"context"
	"encoding/json"
	"io"
	"reflect"
	"testing"
	"time"
)

type testAgent struct {
	conn *AgentConn
	// block makes Prompt wait for cancellation instead of finishing.
	block bool
}

func (a *testAgent) Initialize(_ context.Context, req InitializeRequest) (InitializeResponse, error) {
	return InitializeResponse{ProtocolVersion: req.ProtocolVersion}, nil
}

func (a *testAgent) Authenticate(context.Context, AuthenticateRequest) (AuthenticateResponse, error) {
	return AuthenticateResponse{}, nil
}

func (a *testAgent) NewSession(context.Context, NewSessionRequest) (NewSessionResponse, error) {
	return NewSessionResponse{SessionID: "s1"}, nil
}

func (a *testAgent) Prompt(ctx context.Context, req PromptRequest) (PromptResponse, error) {
	if a.block {
		<-ctx.Done()
		return PromptResponse{}, ctx.Err()
	}
	chunk := SessionUpdate{AgentMessageChunk: &ContentChunk{Content: ContentBlock{Text: &TextContent{Text: "hi"}}}}
	if err := a.conn.SessionUpdate(ctx, SessionNotification{SessionID: req.SessionID, Update: chunk}); err != nil {
		return PromptResponse{}, err
	}
	perm, err := a.conn.RequestPermission(ctx, RequestPermissionRequest{
		SessionID: req.SessionID,
		ToolCall:  ToolCallUpdate{ToolCallID: "t1"},
		Options:   []PermissionOption{{OptionID: "yes", Name: "Allow", Kind: PermissionOptionKindAllowOnce}},
	})
	if err != nil {
		return PromptResponse{}, err
	}
	if perm.Outcome.Selected == nil || perm.Outcome.Selected.OptionID != "yes" {
		return PromptResponse{StopReason: StopReasonRefusal}, nil
	}
	return PromptResponse{StopReason: StopReasonEndTurn}, nil
}

func (a *testAgent) Cancel(context.Context, CancelNotification) error { return nil }

type testClient struct {
	updates chan SessionUpdate
}

func (c *testClient) SessionUpdate(_ context.Context, n SessionNotification) error {
	c.updates <- n.Update
	return nil
}

func (c *testClient) RequestPermission(context.Context, RequestPermissionRequest) (RequestPermissionResponse, error) {
	return RequestPermissionResponse{Outcome: RequestPermissionOutcome{Selected: &SelectedPermissionOutcome{OptionID: "yes"}}}, nil
}

func pair(t *testing.T, agent *testAgent) (*ClientConn, *testClient) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	toAgentR, toAgentW := io.Pipe()
	toClientR, toClientW := io.Pipe()
	agent.conn = ServeAgent(ctx, agent, toAgentR, toClientW)
	client := &testClient{updates: make(chan SessionUpdate, 8)}
	cc := Connect(ctx, client, toClientR, toAgentW)
	t.Cleanup(func() { _ = cc.Close(); _ = agent.conn.Close() })
	return cc, client
}

func TestPromptStreamsAsksAndEnds(t *testing.T) {
	cc, client := pair(t, &testAgent{})
	ctx := context.Background()
	if init, err := cc.Initialize(ctx, InitializeRequest{ProtocolVersion: V1}); err != nil || init.ProtocolVersion != V1 {
		t.Fatalf("initialize = %+v, %v", init, err)
	}
	sess, err := cc.NewSession(ctx, NewSessionRequest{Cwd: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := cc.Prompt(ctx, PromptRequest{SessionID: sess.SessionID, Prompt: []ContentBlock{{Text: &TextContent{Text: "go"}}}})
	if err != nil || resp.StopReason != StopReasonEndTurn {
		t.Fatalf("prompt = %+v, %v", resp, err)
	}
	if u := <-client.updates; u.AgentMessageChunk == nil || u.AgentMessageChunk.Content.Text.Text != "hi" {
		t.Fatalf("update = %+v", u)
	}
}

// A session/cancel stops the running prompt, which answers "cancelled" rather
// than an error.
func TestCancelEndsPromptAsCancelled(t *testing.T) {
	cc, _ := pair(t, &testAgent{block: true})
	ctx := context.Background()
	done := make(chan PromptResponse, 1)
	go func() {
		resp, err := cc.Prompt(ctx, PromptRequest{SessionID: "s1"})
		if err != nil {
			t.Errorf("prompt: %v", err)
		}
		done <- resp
	}()
	time.Sleep(50 * time.Millisecond) // let the prompt reach the agent
	if err := cc.Cancel(ctx, CancelNotification{SessionID: "s1"}); err != nil {
		t.Fatal(err)
	}
	if resp := <-done; resp.StopReason != StopReasonCancelled {
		t.Fatalf("stop reason = %q, want cancelled", resp.StopReason)
	}
}

func TestUnimplementedOptionalMethodIsNotFound(t *testing.T) {
	cc, _ := pair(t, &testAgent{})
	if _, err := cc.LoadSession(context.Background(), LoadSessionRequest{SessionID: "s1"}); err == nil {
		t.Fatal("session/load on an agent without SessionLoader succeeded")
	}
}

// Every tagged variant survives a decode and re-encode, and an unknown one is
// kept verbatim.
func TestUnionsRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		into any
		json string
	}{
		{new(SessionUpdate), `{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hi"}}`},
		{new(SessionUpdate), `{"sessionUpdate":"agent_thought_chunk","content":{"type":"text","text":"hm"}}`},
		{new(SessionUpdate), `{"sessionUpdate":"tool_call","toolCallId":"t1","title":"Read","kind":"read","status":"pending"}`},
		{new(SessionUpdate), `{"sessionUpdate":"tool_call_update","toolCallId":"t1","status":"completed"}`},
		{new(SessionUpdate), `{"sessionUpdate":"plan","entries":[{"content":"a","priority":"high","status":"pending"}]}`},
		{new(SessionUpdate), `{"sessionUpdate":"current_mode_update","currentModeId":"auto"}`},
		{new(SessionUpdate), `{"sessionUpdate":"brand_new_kind","x":1}`},
		{new(ContentBlock), `{"type":"image","data":"AA==","mimeType":"image/png"}`},
		{new(ToolCallContent), `{"type":"diff","path":"a.go","oldText":"a","newText":"b"}`},
		{new(RequestPermissionOutcome), `{"outcome":"cancelled"}`},
		{new(RequestPermissionOutcome), `{"outcome":"selected","optionId":"yes"}`},
		{new(MCPServer), `{"name":"fs","command":"mcp-fs","args":[],"env":[]}`},
		{new(MCPServer), `{"type":"http","name":"web","url":"https://x","headers":[]}`},
	} {
		if err := json.Unmarshal([]byte(tc.json), tc.into); err != nil {
			t.Fatalf("%s: %v", tc.json, err)
		}
		out, err := json.Marshal(tc.into)
		if err != nil {
			t.Fatalf("%s: marshal: %v", tc.json, err)
		}
		var want, got any
		_ = json.Unmarshal([]byte(tc.json), &want)
		_ = json.Unmarshal(out, &got)
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("round trip\n got %s\nwant %s", out, tc.json)
		}
	}
}
