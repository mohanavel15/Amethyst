package server

import (
	"amethyst/protocol"
	"sync"
)

type ServeMux struct {
	handlers map[protocol.State]map[byte]Handler
	mu       sync.RWMutex
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		handlers: map[protocol.State]map[byte]Handler{
			protocol.StateHandshaking: {},
			protocol.StateStatus:      {},
			protocol.StateLogin:       {},
			protocol.StatePlay:        {},
		},
		mu: sync.RWMutex{},
	}
}

func (mux *ServeMux) Handle(state protocol.State, packetID byte, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if handler == nil {
		panic("plasma: nil handler")
	}

	mux.handlers[state][packetID] = handler
}

func (mux *ServeMux) HandleFunc(state protocol.State, packetID byte, handler func(ctx *Context)) {
	if handler == nil {
		panic("plasma: nil handler")
	}

	mux.Handle(state, packetID, HandlerFunc(handler))
}

func (mux *ServeMux) Handler(ctx *Context) (Handler, byte) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	handler, ok := mux.handlers[ctx.conn.state][ctx.Packet.ID]
	if !ok {
		return nil, ctx.Packet.ID
	}

	return handler, ctx.Packet.ID
}

func (mux *ServeMux) ServeProtocol(ctx *Context) {
	handler, _ := mux.Handler(ctx)
	if handler == nil {
		return
	}

	handler.ServeProtocol(ctx)
}

var DefaultServeMux = NewServeMux()

func Handle(state protocol.State, packetID byte, handler Handler) {
	DefaultServeMux.Handle(state, packetID, handler)
}

func HandleFunc(state protocol.State, packetID byte, handler func(ctx *Context)) {
	DefaultServeMux.HandleFunc(state, packetID, handler)
}
