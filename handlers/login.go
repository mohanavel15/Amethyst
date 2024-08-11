package handlers

import (
	"amethyst/protocol"
	"amethyst/protocol/packets/login"
	"amethyst/protocol/packets/play"
	"amethyst/server"
	"log"
)

func LoginStart(w server.ResponseWriter, r *server.Request) {
	loginPack, err := login.UnmarshalServerBoundLoginStart(r.Packet)
	if err != nil {
		w.WritePacket(login.ClientBoundDisconnect{
			Reason: "Invalid Packet",
		}.Marshal())
	}

	server := r.Server()
	server.AddPlayer(r, string(loginPack.Name))

	vt, err := server.SessionEncrypter.GenerateVerifyToken(r)
	if err != nil {
		w.WritePacket(login.ClientBoundDisconnect{
			Reason: "Unable To Connect",
		}.Marshal())
	}

	er := login.ClientBoundEncryptionRequest{
		ServerID:    protocol.String(server.ID),
		PublicKey:   server.SessionEncrypter.PublicKey(),
		VerifyToken: vt,
	}

	w.WritePacket(er.Marshal())
}

func EncryptionResponse(w server.ResponseWriter, r *server.Request) {
	encRes, err := login.UnmarshalServerBoundEncryptionResponse(r.Packet)
	if err != nil {
		w.WritePacket(login.ClientBoundDisconnect{
			Reason: "Invalid Packet",
		}.Marshal())
	}

	server := r.Server()
	sharedSecret, err := server.SessionEncrypter.DecryptAndVerifySharedSecret(r, encRes.SharedSecret, encRes.VerifyToken)
	if err != nil {
		w.WritePacket(login.ClientBoundDisconnect{
			Reason: "Unable To Connect",
		}.Marshal())
	}

	w.SetEncryption(sharedSecret)

	player := server.Player(r)
	log.Println(player.Username())

	setCompress := login.ClientBoundSetCompression{
		Threshold: 1024,
	}
	w.WritePacket(setCompress.Marshal())
	w.SetCompression(1024)

	loginSuccess := login.ClientBoundLoginSuccess{
		UUID:     protocol.String(player.UUID().String()),
		Username: protocol.String(player.Username()),
	}

	w.WritePacket(loginSuccess.Marshal())
	w.SetState(protocol.StatePlay)

	joinGame := play.ClientBoundJoinGame{
		EntityID:         1,
		Gamemode:         play.GamemodeSurvival,
		Dimension:        play.DimensionOverworld,
		Difficulty:       play.DifficultyNormal,
		MaxPlayers:       protocol.UnsignedByte(server.MaxPlayers),
		LevelType:        "default",
		ReducedDebugInfo: false,
	}

	w.WritePacket(joinGame.Marshal())

	spawnPosition := play.ClientBoundSpawnPosition{
		Location: protocol.Position{
			X: 0,
			Y: 65,
			Z: 0,
		},
	}

	w.WritePacket(spawnPosition.Marshal())

	playerPos := play.ClientBoundPlayerPositionAndLook{
		X:     0,
		Y:     65,
		Z:     0,
		Yaw:   0,
		Pitch: 0,
		Flags: 0,
	}
	w.WritePacket(playerPos.Marshal())
}
