package play

import (
	"amethyst/protocol"
)

const ClientBoundPlayerPositionAndLookPacketID byte = 0x08

type ClientBoundPlayerPositionAndLook struct {
	X     protocol.Double
	Y     protocol.Double
	Z     protocol.Double
	Yaw   protocol.Float
	Pitch protocol.Float
	Flags protocol.Byte
}

func (pk ClientBoundPlayerPositionAndLook) Marshal() protocol.Packet {
	return protocol.MarshalPacket(
		ClientBoundPlayerPositionAndLookPacketID,
		pk.X,
		pk.Y,
		pk.Z,
		pk.Yaw,
		pk.Pitch,
		pk.Flags,
	)
}

func UnmarshalClientBoundPlayerPositionAndLook(packet protocol.Packet) (ClientBoundPlayerPositionAndLook, error) {
	var pk ClientBoundPlayerPositionAndLook

	if packet.ID != ClientBoundPlayerPositionAndLookPacketID {
		return pk, protocol.ErrInvalidPacketID
	}

	if err := packet.Scan(
		&pk.X,
		&pk.Y,
		&pk.Z,
		&pk.Yaw,
		&pk.Pitch,
		&pk.Flags,
	); err != nil {
		return pk, err
	}

	return pk, nil
}
