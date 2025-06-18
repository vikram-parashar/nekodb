package main

import (
	"strings"
)

func (s *Server) ExecuteCmd(cmd string, args []DataType) []byte {
	cmd = strings.ToUpper(cmd)

	switch cmd {
	case "PING":
		return s.PingCmd()
	case "SET":
		return s.SetCmd(args)
	case "GET":
		return s.GetCmd(args)
	default:
		return []byte("-ERR unknown command\r\n")
	}
}

func (s *Server) PingCmd() []byte {
	return ([]byte("+PONG\r\n"))
}

func (s *Server) SetCmd(args []DataType) []byte {
	if len(args) != 2 {
		return []byte("-ERR wrong number of arguments for 'set' command\r\n")
	}
	key, val := args[0], args[1]
	if key.Name != "bulk" {
		return []byte("-ERR first argument must be a bulk string\r\n")
	}
	s.kvMu.Lock()
	defer s.kvMu.Unlock()
	
	s.kv.Strings[key.bulk] = val.bulk
	return []byte("+OK\r\n")
}

func (s *Server) GetCmd(args []DataType) []byte {
	if len(args) != 1 {
		return []byte("-ERR wrong number of arguments for 'get' command\r\n")
	}
	key := args[0]
	if key.Name != "bulk" {
		return []byte("-ERR first argument must be a bulk string\r\n")
	}

	s.kvMu.RLock()
	defer s.kvMu.RUnlock()

	if val, ok := s.kv.Strings[key.bulk]; !ok {
		return []byte("-ERR key not found\r\n")
	} else {
		return []byte("+" + val + "\r\n")
	}
}
