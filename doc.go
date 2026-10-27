// Package acp implements the Agent Client Protocol (ACP) v1, the JSON-RPC
// protocol between editors or other clients and coding agents
// (https://agentclientprotocol.com).
//
// An agent implements [Agent] and serves it with [NewAgentConn] and
// [AgentConn.Serve]; a client implements [Client] and reaches an agent with
// [Spawn] or [Connect]. The types follow the stable v1 schema vendored under
// schema/v1.
package acp
