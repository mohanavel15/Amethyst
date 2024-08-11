package handlers

import (
	"amethyst/server"
)

func KeepAlive(ctx *server.Context) {
	conn := ctx.Conn()
	conn.UpdateKeepAlive()
}
