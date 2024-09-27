package play

import (
	"amethyst/protocol"
)

const ClientBoundSpawnPlayerPacketID byte = 0x0C

type ClientBoundSpawnPlayer struct {
	EntityID    protocol.VarInt
	PlayerUUID  protocol.UUID
	X           protocol.Double
	Y           protocol.Double
	Z           protocol.Double
	Yaw         protocol.Angle
	Pitch       protocol.Angle
	CurrentItem protocol.Short
	MetaData    protocol.EntityMetaData
}

func (pk ClientBoundSpawnPlayer) Marshal() protocol.Packet {
	return protocol.MarshalPacket(
		ClientBoundSpawnPlayerPacketID,
		pk.EntityID,
		pk.PlayerUUID,
		pk.X,
		pk.Y,
		pk.Z,
		pk.Yaw,
		pk.Pitch,
		pk.CurrentItem,
		pk.MetaData,
	)
}

func UnmarshalClientBoundSpawnPlayer(packet protocol.Packet) (ClientBoundSpawnPlayer, error) {
	var pk ClientBoundSpawnPlayer

	if packet.ID != ClientBoundSpawnPlayerPacketID {
		return pk, protocol.ErrInvalidPacketID
	}

	if err := packet.Scan(
		&pk.EntityID,
		&pk.PlayerUUID,
		&pk.X,
		&pk.Y,
		&pk.Z,
		&pk.Yaw,
		&pk.Pitch,
		&pk.CurrentItem,
		// &pk.MetaData, Fix...
	); err != nil {
		return pk, err
	}

	return pk, nil
}
