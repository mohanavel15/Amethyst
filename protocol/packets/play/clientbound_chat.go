package play

import "amethyst/protocol"

const ClientBoundChatPacketID = 0x02

type ClientBoundChat struct {
	Message  protocol.Chat
	Position protocol.Byte
}

func (pk ClientBoundChat) Marshal() protocol.Packet {
	return protocol.MarshalPacket(
		ClientBoundChatPacketID,
		pk.Message,
		pk.Position,
	)
}

func UnmarshalClientBoundChat(packet protocol.Packet) (ClientBoundChat, error) {
	var pk ClientBoundChat

	if packet.ID != ClientBoundChatPacketID {
		return pk, protocol.ErrInvalidPacketID
	}

	if err := packet.Scan(
		&pk.Message,
		&pk.Position,
	); err != nil {
		return pk, err
	}

	return pk, nil
}
