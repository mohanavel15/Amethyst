package play

import (
	"amethyst/protocol"
)

const (
	ClientBoundSpawnPositionPacketID byte = 0x05
)

type ClientBoundSpawnPosition struct {
	Location protocol.Position
}

func (pk ClientBoundSpawnPosition) Marshal() protocol.Packet {
	return protocol.MarshalPacket(
		ClientBoundSpawnPositionPacketID,
		pk.Location,
	)
}

func UnmarshalClientBoundSpawnPosition(packet protocol.Packet) (ClientBoundSpawnPosition, error) {
	var pk ClientBoundSpawnPosition

	if packet.ID != ClientBoundSpawnPositionPacketID {
		return pk, protocol.ErrInvalidPacketID
	}

	if err := packet.Scan(
		&pk.Location,
	); err != nil {
		return pk, err
	}

	return pk, nil
}
