package handlers

import "amethyst/server"

func Echo(w server.ResponseWriter, r *server.Request) {
	w.WritePacket(r.Packet)
}
