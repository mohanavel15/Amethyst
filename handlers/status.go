package handlers

import (
	"amethyst/protocol"
	"amethyst/protocol/packets/status"
	"amethyst/server"
)

func StatusRequest(ctx *server.Context) {
	statusResponse := server.StatusResponse{
		Version: server.Version{
			Name:           "Amethyst 1.8.9",
			ProtocolNumber: 47,
		},
		PlayersInfo: ctx.Server().PlayersInfo(),
		IconPath:    "",
		MOTD:        "Amethyst 1.8.9",
	}

	buffer, err := statusResponse.JSON()
	if err != nil {
		return
	}

	ctx.WritePacket(status.ClientBoundResponse{
		JSONResponse: protocol.String(buffer),
	}.Marshal())
}
