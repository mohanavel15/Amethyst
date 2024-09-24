package play

import "amethyst/protocol"

const (
	ActionAddPlayer    = protocol.VarInt(0)
	ActionRemovePlayer = protocol.VarInt(4)
)

const ClientBoundPlayerTabListItemPacketID = 0x38

type ClientBoundPlayerTabListItem struct {
	Action      protocol.VarInt
	PlayerCount protocol.VarInt
	ActionData  protocol.FieldEncoder
}

func (pk ClientBoundPlayerTabListItem) Marshal() protocol.Packet {
	return protocol.MarshalPacket(
		ClientBoundPlayerTabListItemPacketID,
		pk.Action,
		pk.PlayerCount,
		pk.ActionData,
	)
}

type AddPlayerProperty struct {
	Name  protocol.String
	Value protocol.String
}

func (pr AddPlayerProperty) Encode() []byte {
	buf := []byte{}
	buf = append(buf, pr.Name.Encode()...)
	buf = append(buf, pr.Value.Encode()...)
	return buf
}

type TabListActionAddPlayer struct {
	UUID           protocol.UUID
	Name           protocol.String
	Properties     []AddPlayerProperty
	Gamemode       protocol.VarInt
	Ping           protocol.VarInt
	HasDisplayName protocol.Boolean
	DisplayName    protocol.Chat
}

func (pk TabListActionAddPlayer) Encode() []byte {
	buf := []byte{}
	buf = append(buf, pk.UUID.Encode()...)
	buf = append(buf, pk.Name.Encode()...)
	buf = append(buf, protocol.VarInt(len(pk.Properties)).Encode()...)
	for _, prop := range pk.Properties {
		buf = append(buf, prop.Encode()...)
		buf = append(buf, protocol.Boolean(false).Encode()...)
	}
	buf = append(buf, pk.Gamemode.Encode()...)
	buf = append(buf, pk.Ping.Encode()...)
	buf = append(buf, pk.HasDisplayName.Encode()...)
	if pk.HasDisplayName {
		buf = append(buf, pk.DisplayName.Encode()...)
	}
	return buf
}

type TabListActionRemovePlayer struct {
	UUID protocol.UUID
}

func (pk TabListActionRemovePlayer) Encode() []byte {
	buf := []byte{}
	buf = append(buf, pk.UUID.Encode()...)
	return buf
}
