package play

import "amethyst/protocol"

const ServerBoundPlayerPositionPacketID = 0x04

type ServerBoundPlayerPosition struct {
	X        protocol.Double
	Y        protocol.Double
	Z        protocol.Double
	OnGround protocol.Boolean
}

func (pk ServerBoundPlayerPosition) Marshal() protocol.Packet {
	return protocol.MarshalPacket(
		ServerBoundPlayerPositionPacketID,
		pk.X,
		pk.Y,
		pk.Z,
		pk.OnGround,
	)
}

func UnmarshalServerBoundPlayerPosition(packet protocol.Packet) (ServerBoundPlayerPosition, error) {
	var pk ServerBoundPlayerPosition

	if packet.ID != ServerBoundPlayerPositionPacketID {
		return pk, protocol.ErrInvalidPacketID
	}

	if err := packet.Scan(
		&pk.X,
		&pk.Y,
		&pk.Z,
		&pk.OnGround,
	); err != nil {
		return pk, err
	}

	return pk, nil
}
