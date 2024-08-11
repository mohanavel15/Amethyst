package handlers

import "amethyst/server"

func Echo(ctx *server.Context) {
	ctx.WritePacket(ctx.Packet)
}
