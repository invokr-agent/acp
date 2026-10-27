// Command echoagent is the smallest useful ACP agent: it streams each prompt's
// text back as the reply. Point any ACP client (Zed, a messaging bridge) at it
// to check the wiring.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/invokr-ai/acp"
)

type echo struct {
	conn     *acp.AgentConn
	sessions atomic.Int64
}

func (e *echo) Initialize(_ context.Context, _ acp.InitializeRequest) (acp.InitializeResponse, error) {
	return acp.InitializeResponse{ProtocolVersion: acp.V1, AgentInfo: &acp.Implementation{Name: "echoagent", Version: "0.1.0"}}, nil
}

func (e *echo) Authenticate(context.Context, acp.AuthenticateRequest) (acp.AuthenticateResponse, error) {
	return acp.AuthenticateResponse{}, nil
}

func (e *echo) NewSession(context.Context, acp.NewSessionRequest) (acp.NewSessionResponse, error) {
	return acp.NewSessionResponse{SessionID: acp.SessionID(fmt.Sprint("echo-", e.sessions.Add(1)))}, nil
}

func (e *echo) Prompt(ctx context.Context, req acp.PromptRequest) (acp.PromptResponse, error) {
	for _, block := range req.Prompt {
		if block.Text == nil {
			continue
		}
		update := acp.SessionUpdate{AgentMessageChunk: &acp.ContentChunk{Content: block}}
		if err := e.conn.SessionUpdate(ctx, acp.SessionNotification{SessionID: req.SessionID, Update: update}); err != nil {
			return acp.PromptResponse{}, err
		}
	}
	return acp.PromptResponse{StopReason: acp.StopReasonEndTurn}, nil
}

func (e *echo) Cancel(context.Context, acp.CancelNotification) error { return nil }

func main() {
	agent := &echo{conn: acp.NewAgentConn(os.Stdin, os.Stdout)}
	agent.conn.Serve(context.Background(), agent)
	<-agent.conn.Done()
	if err := agent.conn.Err(); err != nil {
		log.Fatal(err)
	}
}
