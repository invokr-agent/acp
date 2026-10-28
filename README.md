# acp

Go types and connections for the [Agent Client Protocol](https://agentclientprotocol.com) (ACP) v1: the JSON-RPC protocol between editors or other clients and coding agents.

The package holds the protocol only and does not depend on any particular agent or client. JSON-RPC 2.0 comes from [`go.lsp.dev/jsonrpc2`](https://github.com/go-language-server/jsonrpc2), with ACP's newline-delimited framing.

```
go get github.com/invokr-agent/acp
```

## Coverage

The stable v1 schema ([`schema-v1.26.0`](https://github.com/agentclientprotocol/agent-client-protocol)):

- **Agent methods:** `initialize`, `authenticate`, `logout`, `session/new|load|resume|close|list|delete|prompt|cancel|set_mode|set_config_option`.
- **Client methods:** `session/update`, `session/request_permission`, `fs/read_text_file|write_text_file`, `terminal/create|output|wait_for_exit|kill|release`, `elicitation/create|complete`.
- **Protocol:** `$/cancel_request`.
- **Extensions:** `_`-prefixed methods via the `Extension` interface.

Unions such as `SessionUpdate` and `ContentBlock` have one pointer field per variant. A variant this version does not know is kept verbatim in `Unknown`.

The unstable methods (`providers/*`, `mcp/message`, `session/fork`) are not included. ACP v2 will be added as a sibling package once it is stable.

## Schema

`schema/v1` and `schema/v2` are copies of the official schemas. `make schema` refreshes them from the spec repo (`SCHEMA_REF=<tag>` pins a release; the default is `main`).

`types.go` is maintained by hand. `TestTypesFollowSchema` fails when it drifts from `schema/v1`, naming what is out of line:
- a missing type;
- object fields that differ;
- union tags or enum values that differ;
- a method from `meta.json` that is not wired into a connection.

After `make schema`, run `make test` and bring `types.go` in line.

## Agent

Implement `acp.Agent`. Optional methods (`SessionLoader`, `SessionModeSetter`, and so on) are separate interfaces; implement the ones you support and advertise them in `InitializeResponse`.

```go
agent := &myAgent{conn: acp.NewAgentConn(os.Stdin, os.Stdout)}
agent.conn.Serve(ctx, agent) // the agent holds its conn before any message arrives
<-agent.conn.Done()
```

Inside `Prompt`, stream progress with `conn.SessionUpdate`, and ask the user with `conn.RequestPermission`. When the client sends `session/cancel`, the prompt's context is cancelled. Return the context's error and the client receives the `cancelled` stop reason, as the spec requires.

See [`examples/agent`](examples/agent/main.go) and [`examples/client`](examples/client/main.go).

## Client

Implement `acp.Client` (updates and permission requests), plus `FileSystem`, `Terminals` or `Elicitor` if you offer them.

```go
conn, err := acp.Spawn(ctx, client, exec.Command("some-agent", "acp"))
init, err := conn.Initialize(ctx, acp.InitializeRequest{ProtocolVersion: acp.V1})
sess, err := conn.NewSession(ctx, acp.NewSessionRequest{Cwd: dir})
resp, err := conn.Prompt(ctx, acp.PromptRequest{SessionID: sess.SessionID, Prompt: blocks})
```

`conn.Cancel` notifies the agent and resolves your pending permission prompts for that session as `cancelled`.

## License

Apache-2.0.
