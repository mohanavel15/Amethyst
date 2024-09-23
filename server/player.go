package server

import (
	"github.com/gofrs/uuid"
)

type Skin struct {
	Value     string
	Signature string
}

type player struct {
	*conn
	uuid     uuid.UUID
	username string
	skin     Skin
}

func (p player) UUID() uuid.UUID {
	return p.uuid
}

func (p player) IntUUID() int32 {
	bs := p.uuid.Bytes()
	i := int32(bs[0])<<24 | int32(bs[1])<<16 | int32(bs[2])<<8 | int32(bs[3])
	return i
}

func (p player) Username() string {
	return p.username
}

func (p player) Skin() Skin {
	return p.skin
}

type Player interface {
	Conn

	// UUID returns the uuid of the player
	UUID() uuid.UUID
	IntUUID() int32

	// Username returns the username of the player
	Username() string

	// Skin returns the Skin of the player
	Skin() Skin
}
