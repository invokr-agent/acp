package acp

import "go.lsp.dev/jsonrpc2"

// NewError returns a JSON-RPC error carrying an ACP error code, e.g.
// NewError(ErrAuthRequired, "log in first").
func NewError(code ErrorCode, message string) error {
	return jsonrpc2.NewError(jsonrpc2.Code(code), message)
}
