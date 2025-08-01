package server

import (
	"container/list"
	"nekodb/internal/parser"
	"net"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type Server struct {
	addr string
	ln   net.Listener

	connsMu sync.RWMutex
	conns   map[*net.Conn]struct{}

	kvMu sync.RWMutex
	kv   *KVStore
}
type KVStore struct {
	Strings     map[string]string
	Expirations map[string]time.Time
	Lists       map[string]*list.List
	Sets        map[string]map[string]struct{}
}

func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		ln:      nil,
		connsMu: sync.RWMutex{},
		conns:   make(map[*net.Conn]struct{}),
		kvMu:    sync.RWMutex{},
		kv: &KVStore{
			Strings:     map[string]string{},
			Expirations: map[string]time.Time{},
			Lists:       map[string]*list.List{},
			Sets:        map[string]map[string]struct{}{},
		},
	}
}

func (s *Server) clearExpiryRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			log.Info("Clearing expired keys...")
			s.kvMu.Lock()
			for key, expiration := range s.kv.Expirations {
				if time.Now().After(expiration) {
					delete(s.kv.Strings, key)
					delete(s.kv.Expirations, key)
				}
			}
			s.kvMu.Unlock()
		}
	}
}

func (s *Server) Start() {
	var err error
	s.ln, err = net.Listen("tcp", s.addr)
	if err != nil {
		log.Error("Server failed to start", "err", err)
		os.Exit(0)
	}

	log.Info("Server started", "at addr", s.addr)

	go s.clearExpiryRoutine()

	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		// wait till we get new conn
		conn, err := s.ln.Accept()
		if err != nil {
			log.Error("Failed to accept connetion: ", err)
			continue
		}
		log.Info("New client connected", "addr", conn.RemoteAddr().String())

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	s.connsMu.Lock()
	defer s.connsMu.Unlock()
	s.conns[&conn] = struct{}{}

	if err := s.readLoop(conn); err != nil {
		log.Error("Error reading from connection: ", err)
	}
}

func (s *Server) readLoop(conn net.Conn) error {
	for {
		rp := parser.NewParser(conn)
		data, err := rp.Read()
		if err != nil {
			return err
		}
		if data.Name != "array" {
			log.Error("Invalid data received", "data", data)
			continue
		}
		if len(data.Arr) == 0 {
			log.Error("Empty array received")
			continue
		}

		cmd := data.Arr[0].Bulk
		args := data.Arr[1:]
		res := s.ExecuteCmd(cmd, args)
		conn.Write(res)
	}
}
