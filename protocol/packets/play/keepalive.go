package play

import "amethyst/protocol"

const KeepAlivePacketID = 0x00

type KeepAlive struct {
	ID protocol.VarInt
}

func (pk KeepAlive) Marshal() protocol.Packet {
	return protocol.MarshalPacket(
		KeepAlivePacketID,
		pk.ID,
	)
}

func UnmarshalKeepAlive(packet protocol.Packet) (KeepAlive, error) {
	var pk KeepAlive

	if packet.ID != KeepAlivePacketID {
		return pk, protocol.ErrInvalidPacketID
	}

	if err := packet.Scan(
		&pk.ID,
	); err != nil {
		return pk, err
	}

	return pk, nil
}
