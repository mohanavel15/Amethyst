package handlers

import (
	"amethyst/protocol"
	"amethyst/protocol/packets/play"
	"amethyst/server"
)

func JoinGame(ctx *server.Context) {
	srv := ctx.Server()
	joinGame := play.ClientBoundJoinGame{
		EntityID:         1,
		Gamemode:         play.GamemodeCreative,
		Dimension:        play.DimensionOverworld,
		Difficulty:       play.DifficultyNormal,
		MaxPlayers:       protocol.UnsignedByte(srv.MaxPlayers),
		LevelType:        "default",
		ReducedDebugInfo: false,
	}

	ctx.WritePacket(joinGame.Marshal())

	spawnPosition := play.ClientBoundSpawnPosition{
		Location: protocol.Position{
			X: 0,
			Y: 65,
			Z: 0,
		},
	}

	ctx.WritePacket(spawnPosition.Marshal())

	playerPos := play.ClientBoundPlayerPositionAndLook{
		X:     0,
		Y:     65,
		Z:     0,
		Yaw:   0,
		Pitch: 0,
		Flags: 0,
	}
	ctx.WritePacket(playerPos.Marshal())
}

func KeepAlive(ctx *server.Context) {
	conn := ctx.Conn()
	conn.UpdateKeepAlive()
}
