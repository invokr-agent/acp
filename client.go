package acp

import (
	"context"
	"errors"
	"io"
	"os/exec"
)

// Client is the baseline every ACP client implements. Optional methods are the
// interfaces below; one the client does not implement answers "method not
// found". Advertise what you support in the InitializeRequest capabilities.
type Client interface {
	// SessionUpdate receives the agent's updates for a session, in order.
	SessionUpdate(ctx context.Context, n SessionNotification) error
	// RequestPermission asks the user. Its ctx is cancelled when this side
	// cancels the session's prompt; returning ctx's error then answers the
	// agent with the "cancelled" outcome, as the spec requires.
	RequestPermission(ctx context.Context, req RequestPermissionRequest) (RequestPermissionResponse, error)
}

type (
	FileSystem interface {
		ReadTextFile(ctx context.Context, req ReadTextFileRequest) (ReadTextFileResponse, error)
		WriteTextFile(ctx context.Context, req WriteTextFileRequest) (WriteTextFileResponse, error)
	}
	Terminals interface {
		CreateTerminal(ctx context.Context, req CreateTerminalRequest) (CreateTerminalResponse, error)
		TerminalOutput(ctx context.Context, req TerminalOutputRequest) (TerminalOutputResponse, error)
		WaitForTerminalExit(ctx context.Context, req WaitForTerminalExitRequest) (WaitForTerminalExitResponse, error)
		KillTerminal(ctx context.Context, req KillTerminalRequest) (KillTerminalResponse, error)
		ReleaseTerminal(ctx context.Context, req ReleaseTerminalRequest) (ReleaseTerminalResponse, error)
	}
	Elicitor interface {
		CreateElicitation(ctx context.Context, req CreateElicitationRequest) (CreateElicitationResponse, error)
		CompleteElicitation(ctx context.Context, n CompleteElicitationNotification) error
	}
)

// ClientConn is the client's end of a connection: it serves the Client and
// calls the agent.
type ClientConn struct {
	*peer
	permissions sessions
	cmd         *exec.Cmd
}

// Connect serves client over r and w, the agent's stdout and stdin.
func Connect(ctx context.Context, client Client, r io.Reader, w io.Writer) *ClientConn {
	c := &ClientConn{}
	m := map[string]method{
		"session/update": notify(client.SessionUpdate),
		"session/request_permission": bind(func(ctx context.Context, req RequestPermissionRequest) (RequestPermissionResponse, error) {
			ctx, done := c.permissions.track(ctx, req.SessionID)
			defer done()
			resp, err := client.RequestPermission(ctx, req)
			if err != nil && ctx.Err() != nil && errors.Is(err, ctx.Err()) {
				return RequestPermissionResponse{Outcome: RequestPermissionOutcome{Cancelled: &RequestPermissionOutcomeCancelled{}}}, nil
			}
			return resp, err
		}),
	}
	if f, ok := client.(FileSystem); ok {
		m["fs/read_text_file"] = bind(f.ReadTextFile)
		m["fs/write_text_file"] = bind(f.WriteTextFile)
	}
	if t, ok := client.(Terminals); ok {
		m["terminal/create"] = bind(t.CreateTerminal)
		m["terminal/output"] = bind(t.TerminalOutput)
		m["terminal/wait_for_exit"] = bind(t.WaitForTerminalExit)
		m["terminal/kill"] = bind(t.KillTerminal)
		m["terminal/release"] = bind(t.ReleaseTerminal)
	}
	if e, ok := client.(Elicitor); ok {
		m["elicitation/create"] = bind(e.CreateElicitation)
		m["elicitation/complete"] = notify(e.CompleteElicitation)
	}
	c.peer = newPeer(ctx, stdio{r, w}, m, client)
	return c
}

// Spawn starts cmd as an agent subprocess and connects client to its stdio.
// The agent's stderr is left as cmd.Stderr; Close also waits for the process.
func Spawn(ctx context.Context, client Client, cmd *exec.Cmd) (*ClientConn, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	c := Connect(ctx, client, stdout, stdin)
	c.cmd = cmd
	return c, nil
}

// Close shuts the connection down and, for a spawned agent, waits for it to exit.
func (c *ClientConn) Close() error {
	err := c.peer.Close()
	if c.cmd != nil {
		err = errors.Join(err, c.cmd.Wait())
	}
	return err
}

func (c *ClientConn) Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error) {
	return call[InitializeResponse](ctx, c.conn, "initialize", req)
}

func (c *ClientConn) Authenticate(ctx context.Context, req AuthenticateRequest) (AuthenticateResponse, error) {
	return call[AuthenticateResponse](ctx, c.conn, "authenticate", req)
}

func (c *ClientConn) Logout(ctx context.Context, req LogoutRequest) (LogoutResponse, error) {
	return call[LogoutResponse](ctx, c.conn, "logout", req)
}

func (c *ClientConn) NewSession(ctx context.Context, req NewSessionRequest) (NewSessionResponse, error) {
	return call[NewSessionResponse](ctx, c.conn, "session/new", req)
}

func (c *ClientConn) LoadSession(ctx context.Context, req LoadSessionRequest) (LoadSessionResponse, error) {
	return call[LoadSessionResponse](ctx, c.conn, "session/load", req)
}

func (c *ClientConn) ResumeSession(ctx context.Context, req ResumeSessionRequest) (ResumeSessionResponse, error) {
	return call[ResumeSessionResponse](ctx, c.conn, "session/resume", req)
}

func (c *ClientConn) CloseSession(ctx context.Context, req CloseSessionRequest) (CloseSessionResponse, error) {
	return call[CloseSessionResponse](ctx, c.conn, "session/close", req)
}

func (c *ClientConn) ListSessions(ctx context.Context, req ListSessionsRequest) (ListSessionsResponse, error) {
	return call[ListSessionsResponse](ctx, c.conn, "session/list", req)
}

func (c *ClientConn) DeleteSession(ctx context.Context, req DeleteSessionRequest) (DeleteSessionResponse, error) {
	return call[DeleteSessionResponse](ctx, c.conn, "session/delete", req)
}

func (c *ClientConn) SetSessionMode(ctx context.Context, req SetSessionModeRequest) (SetSessionModeResponse, error) {
	return call[SetSessionModeResponse](ctx, c.conn, "session/set_mode", req)
}

func (c *ClientConn) SetSessionConfigOption(ctx context.Context, req SetSessionConfigOptionRequest) (SetSessionConfigOptionResponse, error) {
	return call[SetSessionConfigOptionResponse](ctx, c.conn, "session/set_config_option", req)
}

// Prompt runs one turn and blocks until the agent answers with a stop reason;
// updates arrive meanwhile through Client.SessionUpdate.
func (c *ClientConn) Prompt(ctx context.Context, req PromptRequest) (PromptResponse, error) {
	return call[PromptResponse](ctx, c.conn, "session/prompt", req)
}

// Cancel asks the agent to stop the session's turn and resolves this side's
// pending permission requests for it as "cancelled".
func (c *ClientConn) Cancel(ctx context.Context, n CancelNotification) error {
	c.permissions.cancel(n.SessionID)
	return c.conn.Notify(ctx, "session/cancel", n)
}
