package server

import (
	"amethyst/crypto"
	"amethyst/protocol"
	"crypto/aes"

	"github.com/gofrs/uuid"
)

type Context struct {
	Packet protocol.Packet
	server *Server
	conn   *conn
}

func (ctx *Context) SetState(state protocol.State) {
	ctx.conn.state = state
}

func (ctx *Context) WritePacket(p protocol.Packet) error {
	return ctx.conn.WritePacket(p)
}

func (ctx *Context) SetEncryption(sharedSecret []byte) error {
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return err
	}

	ctx.conn.SetCipher(
		crypto.NewEncrypter(block, sharedSecret),
		crypto.NewDecrypter(block, sharedSecret),
	)
	return nil
}

func (ctx *Context) SetCompression(threshold int) {
	ctx.conn.threshold = threshold
}

func (ctx Context) Server() *Server {
	return ctx.server
}

func (ctx Context) ClonePacket() protocol.Packet {
	data := make([]byte, len(ctx.Packet.Data))
	copy(data, ctx.Packet.Data)
	return protocol.Packet{
		ID:   ctx.Packet.ID,
		Data: data,
	}
}

func (ctx Context) ProtocolState() protocol.State {
	return ctx.conn.state
}

func (ctx Context) Conn() Conn {
	return ctx.conn
}

func (ctx *Context) Player() Player {
	return ctx.server.Player(ctx)
}

func (ctx *Context) UpdatePlayerUsername(username string) {
	player := ctx.server.getPlayer(ctx.conn)
	player.username = username
	ctx.server.putPlayer(ctx.conn, player)
}

func (ctx *Context) UpdatePlayerUUID(uuid uuid.UUID) {
	player := ctx.server.getPlayer(ctx.conn)
	player.uuid = uuid
	ctx.server.putPlayer(ctx.conn, player)
}

func (ctx *Context) UpdatePlayerSkin(skin Skin) {
	player := ctx.server.getPlayer(ctx.conn)
	player.skin = skin
	ctx.server.putPlayer(ctx.conn, player)
}
