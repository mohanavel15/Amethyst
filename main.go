package main

import (
	"amethyst/handlers"
	"amethyst/protocol"
	"amethyst/protocol/packets/handshaking"
	"amethyst/protocol/packets/login"
	"amethyst/protocol/packets/status"
	"amethyst/server"
	"log"
)

func main() {
	encrypter, err := server.NewDefaultSessionEncrypter()
	if err != nil {
		log.Panic(err)
	}

	handler := server.NewServeMux()

	handler.HandleFunc(protocol.StateStatus, status.PingPongPacketID, handlers.Echo)
	handler.HandleFunc(protocol.StateStatus, status.ServerBoundRequestPacketID, handlers.StatusRequest)

	handler.HandleFunc(protocol.StateHandshaking, handshaking.ServerBoundHandshakePacketID, handlers.Handshake)

	handler.HandleFunc(protocol.StateLogin, login.ServerBoundLoginStartPacketID, handlers.LoginStart)
	handler.HandleFunc(protocol.StateLogin, login.ServerBoundEncryptionResponsePacketID, handlers.EncryptionResponse)

	srv := &server.Server{
		Addr:                 ":25565",
		Encryption:           true,
		MaxPlayers:           10,
		SessionEncrypter:     encrypter,
		SessionAuthenticator: nil,
		Handler:              handler,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err.Error())
	}
}
