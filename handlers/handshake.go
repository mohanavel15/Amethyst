package handlers

import (
	"amethyst/protocol"
	"amethyst/protocol/packets/handshaking"
	"amethyst/server"
)

func Handshake(ctx *server.Context) {
	hs, err := handshaking.UnmarshalServerBoundHandshake(ctx.Packet)
	if err != nil {
		return
	}

	if hs.IsStatusRequest() {
		ctx.SetState(protocol.StateStatus)
	} else if hs.IsLoginRequest() {
		ctx.SetState(protocol.StateLogin)
	}
}
