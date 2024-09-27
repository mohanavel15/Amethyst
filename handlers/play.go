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

	// https://wiki.vg/index.php?title=Entity_metadata&oldid=7360#Entity
	humanMetaData := protocol.EntityMetaData{}
	// Entity
	humanMetaData.Insert(protocol.EMDByte, 0, 0)
	humanMetaData.Insert(protocol.EMDShort, 1, protocol.Short(0).Encode()...)
	humanMetaData.Insert(protocol.EMDByte, 4, 0)
	// Living Entity
	humanMetaData.Insert(protocol.EMDString, 2, protocol.String(ctx.Player().Username()).Encode()...) // ToDo
	humanMetaData.Insert(protocol.EMDByte, 3, 1)
	humanMetaData.Insert(protocol.EMDFloat, 6, protocol.Float(255).Encode()...)
	humanMetaData.Insert(protocol.EMDInt, 7, protocol.Int(0).Encode()...)
	humanMetaData.Insert(protocol.EMDByte, 8, 0)
	humanMetaData.Insert(protocol.EMDByte, 9, 0)
	humanMetaData.Insert(protocol.EMDByte, 15, 0)
	// Human
	humanMetaData.Insert(protocol.EMDByte, 10, 127)
	humanMetaData.Insert(protocol.EMDByte, 16, 0)
	humanMetaData.Insert(protocol.EMDFloat, 17, protocol.Float(255).Encode()...)
	humanMetaData.Insert(protocol.EMDInt, 18, protocol.Int(0).Encode()...)

	fmt.Println("Reminder: Entity MetaData Not Working...")

	newPlayer := play.ClientBoundSpawnPlayer{
		EntityID:    protocol.VarInt(ctx.Player().IntUUID()),
		PlayerUUID:  protocol.UUID(ctx.Player().UUID()),
		X:           0,
		Y:           65,
		Z:           0,
		Yaw:         0,
		Pitch:       0,
		CurrentItem: 0,
		MetaData:    humanMetaData,
	}

	players := ctx.Server().Players()
	for _, player := range players {
		player.WritePacket(chat.Marshal())
		player.WritePacket(tablist.Marshal())
		if player.IntUUID() != ctx.Player().IntUUID() {
			player.WritePacket(newPlayer.Marshal())
		}
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

func PlayerPosition(ctx *server.Context) {
	pos, err := play.UnmarshalServerBoundPlayerPosition(ctx.Packet)
	if err != nil {
		chat := play.ClientBoundChat{
			Message: protocol.Chat{
				Text:  "Invalid Packet",
				Color: "red",
			},
			Position: play.ChatTypeSystem,
		}

		ctx.WritePacket(chat.Marshal())
		return
	}

	fmt.Printf("| %-10.2f | %-10.2f | %-10.2f | %-10t |\n", pos.X, pos.Y, pos.Z, pos.OnGround)
}
