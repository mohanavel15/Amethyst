package play

import "amethyst/protocol"

const ServerBoundChatPacketID = 0x01

type ServerBoundChat struct {
	Message protocol.String
}

func (pk ServerBoundChat) Marshal() protocol.Packet {
	return protocol.MarshalPacket(
		ServerBoundChatPacketID,
		pk.Message,
	)
}

func UnmarshalServerBoundChat(packet protocol.Packet) (ServerBoundChat, error) {
	var pk ServerBoundChat

	if packet.ID != ServerBoundChatPacketID {
		return pk, protocol.ErrInvalidPacketID
	}

	if err := packet.Scan(
		&pk.Message,
	); err != nil {
		return pk, err
	}

	return pk, nil
}
