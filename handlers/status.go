package handlers

import (
	"amethyst/protocol"
	"amethyst/protocol/packets/status"
	"amethyst/server"
)

func StatusRequest(w server.ResponseWriter, r *server.Request) {
	statusResponse := server.StatusResponse{
		Version: server.Version{
			Name:           "Amethyst 1.8.9",
			ProtocolNumber: 47,
		},
		PlayersInfo: r.Server().PlayersInfo(),
		IconPath:    "",
		MOTD:        "Amethyst 1.8.9",
	}

	bb, err := statusResponse.JSON()
	if err != nil {
		return
	}

	w.WritePacket(status.ClientBoundResponse{
		JSONResponse: protocol.String(bb),
	}.Marshal())
}
