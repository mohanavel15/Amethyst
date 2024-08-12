package handlers

import (
	"amethyst/protocol"
	"amethyst/protocol/packets/login"
	"amethyst/server"
	"log"
)

func LoginStart(ctx *server.Context) {
	loginPack, err := login.UnmarshalServerBoundLoginStart(ctx.Packet)
	if err != nil {
		ctx.WritePacket(login.ClientBoundDisconnect{
			Reason: protocol.Chat{
				Text: "Invalid Packet",
			},
		}.Marshal())
		return
	}

	server := ctx.Server()
	server.AddPlayer(ctx, string(loginPack.Name))

	vt, err := server.SessionEncrypter.GenerateVerifyToken(ctx)
	if err != nil {
		ctx.WritePacket(login.ClientBoundDisconnect{
			Reason: protocol.Chat{
				Text: "Unable To Connect",
			},
		}.Marshal())
		return
	}

	er := login.ClientBoundEncryptionRequest{
		ServerID:    protocol.String(server.ID),
		PublicKey:   server.SessionEncrypter.PublicKey(),
		VerifyToken: vt,
	}

	ctx.WritePacket(er.Marshal())
}

func EncryptionResponse(ctx *server.Context) {
	encRes, err := login.UnmarshalServerBoundEncryptionResponse(ctx.Packet)
	if err != nil {
		ctx.WritePacket(login.ClientBoundDisconnect{
			Reason: protocol.Chat{
				Text: "Invalid Packet",
			},
		}.Marshal())
		return
	}

	server := ctx.Server()
	sharedSecret, err := server.SessionEncrypter.DecryptAndVerifySharedSecret(ctx, encRes.SharedSecret, encRes.VerifyToken)
	if err != nil {
		ctx.WritePacket(login.ClientBoundDisconnect{
			Reason: protocol.Chat{
				Text: "Unable To Connect",
			},
		}.Marshal())
		return
	}

	ctx.SetEncryption(sharedSecret)

	player := server.Player(ctx)
	log.Println(player.Username())

	setCompress := login.ClientBoundSetCompression{
		Threshold: 1024,
	}
	ctx.WritePacket(setCompress.Marshal())
	ctx.SetCompression(1024)

	loginSuccess := login.ClientBoundLoginSuccess{
		UUID:     protocol.String(player.UUID().String()),
		Username: protocol.String(player.Username()),
	}

	ctx.WritePacket(loginSuccess.Marshal())
	ctx.SetState(protocol.StatePlay)

	JoinGame(ctx)
}
