package handlers

import (
	"amethyst/protocol"
	"amethyst/protocol/packets/play"
	"amethyst/server"
	"fmt"
	"time"
)

func JoinGame(ctx *server.Context) {
	time.Sleep(5 * time.Millisecond)
	srv := ctx.Server()
	joinGame := play.ClientBoundJoinGame{
		EntityID:         protocol.Int(ctx.Player().IntUUID()),
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

	formatted := fmt.Sprintf("%s joined the game", ctx.Player().Username())
	chat := play.ClientBoundChat{
		Message: protocol.Chat{
			Text: formatted,
		},
		Position: play.ChatTypeSystem,
	}

	tablist := play.ClientBoundPlayerTabListItem{
		Action:      play.ActionAddPlayer,
		PlayerCount: 1,
		ActionData: play.TabListActionAddPlayer{
			UUID: protocol.UUID(ctx.Player().UUID()),
			Name: protocol.String(ctx.Player().Username()),
			Properties: []play.AddPlayerProperty{
				{
					Name:  "texture",
					Value: "",
				},
			},
			Gamemode:       1,
			Ping:           2,
			HasDisplayName: false,
		},
	}

	players := ctx.Server().Players()
	for _, player := range players {
		player.WritePacket(chat.Marshal())
		player.WritePacket(tablist.Marshal())
	}
}

func KeepAlive(ctx *server.Context) {
	conn := ctx.Conn()
	conn.UpdateKeepAlive()
}

func Chat(ctx *server.Context) {
	chat_string, err := play.UnmarshalServerBoundChat(ctx.Packet)
	if err != nil {
		chat := play.ClientBoundChat{
			Message: protocol.Chat{
				Text:  "Invalid Chat Message",
				Color: "red",
			},
			Position: play.ChatTypeSystem,
		}

		ctx.WritePacket(chat.Marshal())
		return
	}

	formatted := fmt.Sprintf("[%s] %s", ctx.Player().Username(), string(chat_string.Message))
	chat := play.ClientBoundChat{
		Message: protocol.Chat{
			Text: formatted,
		},
		Position: play.ChatTypeChat,
	}

	players := ctx.Server().Players()
	for _, player := range players {
		player.WritePacket(chat.Marshal())
	}
}
