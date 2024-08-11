package play

import (
	"amethyst/protocol"
)

const (
	ClientBoundJoinGamePacketID byte = 0x01
)

type ClientBoundJoinGame struct {
	EntityID         protocol.Int
	Gamemode         protocol.UnsignedByte
	Dimension        protocol.Byte
	Difficulty       protocol.UnsignedByte
	MaxPlayers       protocol.UnsignedByte
	LevelType        protocol.String
	ReducedDebugInfo protocol.Boolean
}

// type DimensionCodecVanilla struct {
// 	Name               string  `nbt:"name"`
// 	PiglinSafe         byte    `nbt:"piglin_safe"`
// 	Natural            byte    `nbt:"natural"`
// 	AmbientLight       float32 `nbt:"ambient_light"`
// 	FixedTime          int     `nbt:"fixed_time"`
// 	Infiniburn         string  `nbt:"infiniburn"`
// 	RespawnAnchorWorks byte    `nbt:"respawn_anchor_works"`
// 	HasSkylight        byte    `nbt:"has_skylight"`
// 	BedWorks           byte    `nbt:"bed_works"`
// 	Effects            string  `nbt:"effects"`
// 	HasRaids           byte    `nbt:"has_raids"`
// 	LogicalHeight      int     `nbt:"logical_height"`
// 	CoordinateScale    float32 `nbt:"coordinate_scale"`
// 	Ultrawarm          byte    `nbt:"ultrawarm"`
// 	HasCeiling         byte    `nbt:"has_ceiling"`
// }

func (pk ClientBoundJoinGame) Marshal() protocol.Packet {
	return protocol.MarshalPacket(
		ClientBoundJoinGamePacketID,
		pk.EntityID,
		pk.Gamemode,
		pk.Dimension,
		pk.Difficulty,
		pk.MaxPlayers,
		pk.LevelType,
		pk.ReducedDebugInfo,
	)
}

func UnmarshalClientBoundJoinGame(packet protocol.Packet) (ClientBoundJoinGame, error) {
	var pk ClientBoundJoinGame

	if packet.ID != ClientBoundJoinGamePacketID {
		return pk, protocol.ErrInvalidPacketID
	}

	if err := packet.Scan(
		&pk.EntityID,
		&pk.Gamemode,
		&pk.Dimension,
		&pk.Difficulty,
		&pk.MaxPlayers,
		&pk.LevelType,
		&pk.ReducedDebugInfo,
	); err != nil {
		return pk, err
	}

	return pk, nil
}
