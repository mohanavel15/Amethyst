package handlers

import (
	"amethyst/protocol"
	"amethyst/protocol/packets/handshaking"
	"amethyst/server"
)

func Handshake(w server.ResponseWriter, r *server.Request) {
	hs, err := handshaking.UnmarshalServerBoundHandshake(r.Packet)
	if err != nil {
		return
	}

	if hs.IsStatusRequest() {
		w.SetState(protocol.StateStatus)
	} else if hs.IsLoginRequest() {
		w.SetState(protocol.StateLogin)
	}
}
