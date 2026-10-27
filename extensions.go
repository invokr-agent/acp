package acp

import (
	"context"
	"encoding/json"
)

// Extension handles "_"-prefixed methods and notifications, ACP's extension
// space. Implement it on an Agent or Client to receive them; without it,
// extension calls get "method not found" and notifications are dropped.
type Extension interface {
	HandleExtension(ctx context.Context, method string, params json.RawMessage) (any, error)
}

// CallExtension sends a "_"-prefixed request and decodes its result into result.
func (p *peer) CallExtension(ctx context.Context, method string, params, result any) error {
	_, err := p.conn.Call(ctx, method, params, result)
	return err
}

// NotifyExtension sends a "_"-prefixed notification.
func (p *peer) NotifyExtension(ctx context.Context, method string, params any) error {
	return p.conn.Notify(ctx, method, params)
}
