package acp

import (
	"context"
	"errors"
	"io"
)

// Agent is the baseline every ACP agent implements. Optional methods are the
// interfaces below; a method the agent does not implement answers "method not
// found". Advertise what you support in the InitializeResponse capabilities.
type Agent interface {
	Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error)
	Authenticate(ctx context.Context, req AuthenticateRequest) (AuthenticateResponse, error)
	NewSession(ctx context.Context, req NewSessionRequest) (NewSessionResponse, error)
	// Prompt runs one turn. Its ctx is cancelled when the client sends
	// session/cancel for the session; returning ctx's error then becomes the
	// "cancelled" stop reason, as the spec requires.
	Prompt(ctx context.Context, req PromptRequest) (PromptResponse, error)
	// Cancel is told about a session/cancel, after the prompt's ctx is cancelled.
	Cancel(ctx context.Context, n CancelNotification) error
}

type (
	SessionLoader interface {
		LoadSession(ctx context.Context, req LoadSessionRequest) (LoadSessionResponse, error)
	}
	SessionResumer interface {
		ResumeSession(ctx context.Context, req ResumeSessionRequest) (ResumeSessionResponse, error)
	}
	SessionCloser interface {
		CloseSession(ctx context.Context, req CloseSessionRequest) (CloseSessionResponse, error)
	}
	SessionLister interface {
		ListSessions(ctx context.Context, req ListSessionsRequest) (ListSessionsResponse, error)
	}
	SessionDeleter interface {
		DeleteSession(ctx context.Context, req DeleteSessionRequest) (DeleteSessionResponse, error)
	}
	SessionModeSetter interface {
		SetSessionMode(ctx context.Context, req SetSessionModeRequest) (SetSessionModeResponse, error)
	}
	SessionConfigSetter interface {
		SetSessionConfigOption(ctx context.Context, req SetSessionConfigOptionRequest) (SetSessionConfigOptionResponse, error)
	}
	LogoutHandler interface {
		Logout(ctx context.Context, req LogoutRequest) (LogoutResponse, error)
	}
)

// AgentConn is the agent's end of a connection: it serves the Agent and calls
// the client.
type AgentConn struct {
	*peer
	prompts sessions
}

// NewAgentConn makes the agent's end of a connection over r and w (a client
// usually launches the agent and talks to its stdin and stdout). Hand it to the
// agent, then call Serve: the agent has its connection before any message is
// read.
func NewAgentConn(r io.Reader, w io.Writer) *AgentConn {
	return &AgentConn{peer: newPeer(stdio{r, w})}
}

// Serve dispatches incoming messages to agent until the stream ends or ctx is
// cancelled. It returns at once; wait on Done.
func (c *AgentConn) Serve(ctx context.Context, agent Agent) {
	m := map[string]method{
		"initialize":   bind(agent.Initialize),
		"authenticate": bind(agent.Authenticate),
		"session/new":  bind(agent.NewSession),
		"session/prompt": bind(func(ctx context.Context, req PromptRequest) (PromptResponse, error) {
			ctx, done := c.prompts.track(ctx, req.SessionID)
			defer done()
			resp, err := agent.Prompt(ctx, req)
			if err != nil && ctx.Err() != nil && errors.Is(err, ctx.Err()) {
				return PromptResponse{StopReason: StopReasonCancelled}, nil
			}
			return resp, err
		}),
		"session/cancel": notify(func(ctx context.Context, n CancelNotification) error {
			c.prompts.cancel(n.SessionID)
			return agent.Cancel(ctx, n)
		}),
	}
	opt(m, agent, "session/load", SessionLoader.LoadSession)
	opt(m, agent, "session/resume", SessionResumer.ResumeSession)
	opt(m, agent, "session/close", SessionCloser.CloseSession)
	opt(m, agent, "session/list", SessionLister.ListSessions)
	opt(m, agent, "session/delete", SessionDeleter.DeleteSession)
	opt(m, agent, "session/set_mode", SessionModeSetter.SetSessionMode)
	opt(m, agent, "session/set_config_option", SessionConfigSetter.SetSessionConfigOption)
	opt(m, agent, "logout", LogoutHandler.Logout)
	c.serve(ctx, m, agent)
}

// SessionUpdate streams one update (message chunk, tool call, plan, ...) to the client.
func (c *AgentConn) SessionUpdate(ctx context.Context, n SessionNotification) error {
	return c.conn.Notify(ctx, "session/update", n)
}

func (c *AgentConn) RequestPermission(ctx context.Context, req RequestPermissionRequest) (RequestPermissionResponse, error) {
	return call[RequestPermissionResponse](ctx, c.conn, "session/request_permission", req)
}

func (c *AgentConn) ReadTextFile(ctx context.Context, req ReadTextFileRequest) (ReadTextFileResponse, error) {
	return call[ReadTextFileResponse](ctx, c.conn, "fs/read_text_file", req)
}

func (c *AgentConn) WriteTextFile(ctx context.Context, req WriteTextFileRequest) (WriteTextFileResponse, error) {
	return call[WriteTextFileResponse](ctx, c.conn, "fs/write_text_file", req)
}

func (c *AgentConn) CreateTerminal(ctx context.Context, req CreateTerminalRequest) (CreateTerminalResponse, error) {
	return call[CreateTerminalResponse](ctx, c.conn, "terminal/create", req)
}

func (c *AgentConn) TerminalOutput(ctx context.Context, req TerminalOutputRequest) (TerminalOutputResponse, error) {
	return call[TerminalOutputResponse](ctx, c.conn, "terminal/output", req)
}

func (c *AgentConn) WaitForTerminalExit(ctx context.Context, req WaitForTerminalExitRequest) (WaitForTerminalExitResponse, error) {
	return call[WaitForTerminalExitResponse](ctx, c.conn, "terminal/wait_for_exit", req)
}

func (c *AgentConn) KillTerminal(ctx context.Context, req KillTerminalRequest) (KillTerminalResponse, error) {
	return call[KillTerminalResponse](ctx, c.conn, "terminal/kill", req)
}

func (c *AgentConn) ReleaseTerminal(ctx context.Context, req ReleaseTerminalRequest) (ReleaseTerminalResponse, error) {
	return call[ReleaseTerminalResponse](ctx, c.conn, "terminal/release", req)
}

func (c *AgentConn) CreateElicitation(ctx context.Context, req CreateElicitationRequest) (CreateElicitationResponse, error) {
	return call[CreateElicitationResponse](ctx, c.conn, "elicitation/create", req)
}

func (c *AgentConn) CompleteElicitation(ctx context.Context, n CompleteElicitationNotification) error {
	return c.conn.Notify(ctx, "elicitation/complete", n)
}
