// Command client is the smallest useful ACP client: it spawns an agent, sends
// one prompt and prints the streamed reply, asking on the terminal whenever the
// agent wants permission.
//
//	go run ./examples/client "say hi" go run ./examples/agent
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/invokr-ai/acp"
)

type terminal struct{}

func (terminal) SessionUpdate(_ context.Context, n acp.SessionNotification) error {
	if c := n.Update.AgentMessageChunk; c != nil && c.Content.Text != nil {
		fmt.Print(c.Content.Text.Text)
	}
	return nil
}

func (terminal) RequestPermission(_ context.Context, req acp.RequestPermissionRequest) (acp.RequestPermissionResponse, error) {
	for i, o := range req.Options {
		fmt.Printf("\n[%d] %s", i, o.Name)
	}
	fmt.Print("\nchoice: ")
	var i int
	if _, err := fmt.Fscan(bufio.NewReader(os.Stdin), &i); err != nil || i < 0 || i >= len(req.Options) {
		return acp.RequestPermissionResponse{Outcome: acp.RequestPermissionOutcome{Cancelled: &acp.RequestPermissionOutcomeCancelled{}}}, nil
	}
	return acp.RequestPermissionResponse{Outcome: acp.RequestPermissionOutcome{Selected: &acp.SelectedPermissionOutcome{OptionID: req.Options[i].OptionID}}}, nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: client PROMPT AGENT [ARGS...]")
	}
	ctx := context.Background()
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stderr = os.Stderr
	conn, err := acp.Spawn(ctx, terminal{}, cmd)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	if _, err := conn.Initialize(ctx, acp.InitializeRequest{ProtocolVersion: acp.V1}); err != nil {
		log.Fatal(err)
	}
	wd, _ := os.Getwd()
	sess, err := conn.NewSession(ctx, acp.NewSessionRequest{Cwd: wd, MCPServers: []acp.MCPServer{}})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := conn.Prompt(ctx, acp.PromptRequest{SessionID: sess.SessionID,
		Prompt: []acp.ContentBlock{{Text: &acp.TextContent{Text: strings.TrimSpace(os.Args[1])}}}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n[%s]\n", resp.StopReason)
}
