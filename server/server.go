package server

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/gofrs/uuid"
)

type HandlerFunc func(ctx *Context)

func (f HandlerFunc) ServeProtocol(ctx *Context) {
	f(ctx)
}

type Handler interface {
	ServeProtocol(ctx *Context)
}

type Version struct {
	// Name of the version eg. "1.16.5"
	Name string `json:"name"`
	// Protocol version number eg. 754
	ProtocolNumber int `json:"protocol"`
}

type PlayerInfo struct {
	Name string `json:"name"`
	UUID string `json:"id"`
}

type PlayersInfo struct {
	MaxPlayers    int          `json:"max"`
	PlayersOnline int          `json:"online"`
	Players       []PlayerInfo `json:"sample"`
}

type StatusResponse struct {
	Version     Version
	PlayersInfo PlayersInfo
	IconPath    string
	MOTD        string
}

func (sr StatusResponse) JSON() ([]byte, error) {
	response := struct {
		Version     Version     `json:"version"`
		Players     PlayersInfo `json:"players"`
		Description struct {
			Text string `json:"text"`
		} `json:"description"`
		Favicon string `json:"favicon"`
	}{}

	response.Version = sr.Version
	response.Players = sr.PlayersInfo
	response.Description.Text = sr.MOTD

	if sr.IconPath != "" {
		img64, err := loadImageAndEncodeToBase64String(sr.IconPath)
		if err != nil {
			return nil, err
		}
		response.Favicon = fmt.Sprintf("data:image/png;base64,%s", img64)
	}

	return json.Marshal(response)
}

func loadImageAndEncodeToBase64String(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	imgFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	fileInfo, err := imgFile.Stat()
	if err != nil {
		return "", err
	}

	buffer := make([]byte, fileInfo.Size())
	fileReader := bufio.NewReader(imgFile)
	_, err = fileReader.Read(buffer)
	if err != nil {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(buffer), nil
}

const DefaultAddr string = ":25565"

// Server defines the struct of a running Minecraft server
type Server struct {
	ID         string
	Addr       string
	Encryption bool
	MaxPlayers int

	SessionEncrypter     SessionEncrypter
	SessionAuthenticator SessionAuthenticator
	Handler              Handler

	listener net.Listener
	players  map[*conn]player
	mu       sync.RWMutex
}

func (srv *Server) getPlayer(conn *conn) player {
	srv.mu.RLock()
	defer srv.mu.RUnlock()
	return srv.players[conn]
}

func (srv *Server) putPlayer(c *conn, p player) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.players[c] = p
}

func (srv *Server) removePlayer(c *conn) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	delete(srv.players, c)
}

func (srv *Server) Player(ctx *Context) Player {
	return srv.getPlayer(ctx.conn)
}

func (srv *Server) Players() []Player {
	srv.mu.RLock()
	defer srv.mu.RUnlock()

	var players []Player
	for _, player := range srv.players {
		players = append(players, player)
	}
	return players
}

func (srv *Server) AddPlayer(ctx *Context, username string) {
	srv.putPlayer(ctx.conn, player{
		conn:     ctx.conn,
		uuid:     uuid.NewV3(uuid.NamespaceOID, "OfflinePlayer:"+username),
		username: username,
	})
}

func (srv *Server) PlayersInfo() PlayersInfo {
	pp := srv.Players()
	var players []PlayerInfo

	for _, p := range pp {
		players = append(players, PlayerInfo{
			Name: p.Username(),
			UUID: p.UUID().String(),
		})
	}

	srv.mu.RLock()
	defer srv.mu.RUnlock()
	return PlayersInfo{
		MaxPlayers:    srv.MaxPlayers,
		PlayersOnline: len(pp),
		Players:       players,
	}
}

func (srv *Server) IsRunning() bool {
	return srv.listener != nil
}

func (srv *Server) Close() error {
	if srv.listener != nil {
		return srv.listener.Close()
	}

	return nil
}

func (srv *Server) ListenAndServe() error {
	if srv.listener != nil {
		return errors.New("server is already running")
	}

	addr := srv.Addr
	if addr == "" {
		addr = DefaultAddr
	}

	srv.players = map[*conn]player{}
	srv.mu = sync.RWMutex{}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()
	srv.listener = l

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		go srv.serve(c)
	}
}

func (srv *Server) serve(c net.Conn) {
	conn := wrapConn(c)
	defer func() {
		srv.removePlayer(conn)
		conn.Close()
	}()

	ctx := Context{
		server: srv,
		conn:   conn,
	}

	for {
		pk, err := conn.ReadPacket()
		if err != nil {
			return
		}

		ctx.Packet = pk
		srv.Handler.ServeProtocol(&ctx)
	}
}

func ListenAndServe(addr string, handler Handler) error {
	encrypter, err := NewDefaultSessionEncrypter()
	if err != nil {
		return err
	}

	if handler == nil {
		handler = DefaultServeMux
	}

	srv := &Server{
		Addr:                 addr,
		Encryption:           true,
		SessionEncrypter:     encrypter,
		SessionAuthenticator: &MojangSessionAuthenticator{},
		Handler:              handler,
	}

	return srv.ListenAndServe()
}
