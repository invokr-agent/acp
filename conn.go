package acp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"

	"go.lsp.dev/jsonrpc2"
)

// V1 is the ACP major version this package speaks.
const V1 ProtocolVersion = 1

// NewError returns a JSON-RPC error carrying an ACP error code, e.g.
// NewError(ErrAuthRequired, "log in first").
func NewError(code ErrorCode, message string) error {
	return jsonrpc2.NewError(jsonrpc2.Code(code), message)
}

// Extension handles "_"-prefixed methods and notifications, ACP's extension
// space. Implement it on an Agent or Client to receive them; without it,
// extension calls get "method not found" and notifications are dropped.
type Extension interface {
	HandleExtension(ctx context.Context, method string, params json.RawMessage) (any, error)
}

type method func(ctx context.Context, params json.RawMessage) (any, error)

// bind adapts a typed handler to the dispatch table.
func bind[P, R any](f func(context.Context, P) (R, error)) method {
	return func(ctx context.Context, raw json.RawMessage) (any, error) {
		var p P
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &p); err != nil {
				return nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error())
			}
		}
		return f(ctx, p)
	}
}

// opt registers impl's method under name when impl implements I, the optional
// interface it belongs to.
func opt[I, P, R any](m map[string]method, impl any, name string, f func(I, context.Context, P) (R, error)) {
	if i, ok := impl.(I); ok {
		m[name] = bind(func(ctx context.Context, p P) (R, error) { return f(i, ctx, p) })
	}
}

// notify adapts a typed notification handler.
func notify[P any](f func(context.Context, P) error) method {
	return bind(func(ctx context.Context, p P) (any, error) { return nil, f(ctx, p) })
}

func call[R any](ctx context.Context, c jsonrpc2.Conn, name string, params any) (R, error) {
	var r R
	_, err := c.Call(ctx, name, params, &r)
	return r, err
}

// peer is the side-independent half of a connection: dispatch, $/cancel_request
// and lifecycle. Requests run concurrently (a prompt must be able to call back
// into the peer); notifications run in wire order so updates arrive in order.
type peer struct {
	conn    jsonrpc2.Conn
	methods map[string]method
	ext     Extension
	cancel  func(jsonrpc2.ID)
}

func newPeer(ctx context.Context, rwc io.ReadWriteCloser, methods map[string]method, impl any) *peer {
	p := &peer{methods: methods}
	p.ext, _ = impl.(Extension)
	handler, cancel := jsonrpc2.CancelHandler(p.handle)
	p.cancel = cancel
	p.conn = jsonrpc2.NewConn(jsonrpc2.NewNDJSONStream(rwc), jsonrpc2.WithPreempter(p))
	p.conn.Go(ctx, handler)
	return p
}

func (p *peer) handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	if req.IsCall() {
		jsonrpc2.Async(ctx)
	}
	params := json.RawMessage(req.Params())
	if m, ok := p.methods[req.Method()]; ok {
		return m(ctx, params)
	}
	if strings.HasPrefix(req.Method(), "_") && p.ext != nil {
		return p.ext.HandleExtension(ctx, req.Method(), params)
	}
	if !req.IsCall() {
		return nil, nil // unknown notifications are ignored, per JSON-RPC
	}
	return nil, jsonrpc2.NewError(jsonrpc2.MethodNotFound, "method not found: "+req.Method())
}

// Preempt answers $/cancel_request ahead of the request it cancels.
func (p *peer) Preempt(_ context.Context, req *jsonrpc2.Request) (any, error) {
	if req.Method() != "$/cancel_request" {
		return nil, jsonrpc2.ErrNotHandled
	}
	var n CancelRequestNotification
	if json.Unmarshal(req.Params(), &n) == nil {
		var num int64
		var str string
		switch {
		case json.Unmarshal(n.RequestID, &num) == nil:
			p.cancel(jsonrpc2.NewNumberID(num))
		case json.Unmarshal(n.RequestID, &str) == nil:
			p.cancel(jsonrpc2.NewStringID(str))
		}
	}
	return struct{}{}, nil
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

// Done is closed once the connection has terminated.
func (p *peer) Done() <-chan struct{} { return p.conn.Done() }

// Err reports what terminated the connection; nil for a clean end of stream.
func (p *peer) Err() error { return p.conn.Err() }

// Close shuts the connection down and waits for in-flight work.
func (p *peer) Close() error { return p.conn.Close() }

// sessions tracks one cancel func per session, so a session/cancel can stop the
// work in flight for it.
type sessions struct {
	mu      sync.Mutex
	cancels map[SessionID][]context.CancelFunc
}

func (s *sessions) track(ctx context.Context, id SessionID) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)
	s.mu.Lock()
	if s.cancels == nil {
		s.cancels = map[SessionID][]context.CancelFunc{}
	}
	s.cancels[id] = append(s.cancels[id], cancel)
	n := len(s.cancels[id]) - 1
	s.mu.Unlock()
	return ctx, func() {
		s.mu.Lock()
		if c := s.cancels[id]; n < len(c) {
			c[n] = func() {}
		}
		s.mu.Unlock()
		cancel()
	}
}

func (s *sessions) cancel(id SessionID) {
	s.mu.Lock()
	cancels := s.cancels[id]
	delete(s.cancels, id)
	s.mu.Unlock()
	for _, c := range cancels {
		c()
	}
}

// stdio joins a reader and a writer into the ReadWriteCloser a stream needs.
type stdio struct {
	io.Reader
	io.Writer
}

func (s stdio) Close() error { return errors.Join(closeIf(s.Reader), closeIf(s.Writer)) }

func closeIf(v any) error {
	if c, ok := v.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
