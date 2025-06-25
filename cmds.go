package main

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (s *Server) ExecuteCmd(cmd string, args []DataType) []byte {
	cmd = strings.ToUpper(cmd)

	switch cmd {
	//server commands
	case "PING":
		return s.PingCmd()
	case "ECHO":
		return s.EchoCmd(args)
	//string cmds
	case "SET":
		return s.SetCmd(args)
	case "GET":
		return s.GetCmd(args)
	case "DEL":
		return s.DelCmd(args)
	case "EXISTS":
		return s.ExistsCmd(args)
	//interger cmds
	case "INCR":
		return s.IncrCmd(args)
	case "DECR":
		return s.DecrCmd(args)
	//expiry cmds
	case "EX":
		return s.ExCmd(args)
	case "TTL":
		return s.TtlCmd(args)
	//List cmds
	case "LPUSH":
		return s.LpushCmd(args)
	case "RPUSH":
		return s.RpushCmd(args)
	case "LPOP":
		return s.lpopCmd(args)
	case "RPOP":
		return s.rpopCmd(args)
	default:
		return []byte("-ERR unknown command\r\n")
	}
}

func (s *Server) EchoCmd(args []DataType) []byte {
	if len(args) != 1 {
		return []byte("-ERR wrong number of arguments for 'echo' command\r\n")
	}

	if args[0].Name != "bulk" {
		return []byte("-ERR first argument must be a bulk string\r\n")
	}

	return ([]byte("+" + args[0].bulk + "\r\n"))
}

func (s *Server) PingCmd() []byte {
	return ([]byte("+PONG\r\n"))
}

func (s *Server) SetCmd(args []DataType) []byte {
	if len(args) != 2 {
		return []byte("-ERR wrong number of arguments for 'set' command\r\n")
	}
	key, val := args[0], args[1]
	if key.Name != "bulk" || val.Name != "bulk" {
		return []byte("-ERR first two arguments must be bulk strings\r\n")
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
	key := args[0].bulk
	if key == "" {
		return []byte("-ERR first argument must be a bulk string\r\n")
	}

	s.kvMu.RLock()
	defer s.kvMu.RUnlock()

	s.clearExpiry(key)

	if val, ok := s.kv.Strings[key]; !ok {
		return []byte("-ERR key not found\r\n")
	} else {
		return []byte("+" + val + "\r\n")
	}
}

func (s *Server) DelCmd(args []DataType) []byte {
	if len(args) == 0 {
		return []byte("-ERR no key provided\r\n")
	}
	for _, key := range args {
		if key.Name != "bulk" {
			return []byte("-ERR all arguments must be bulk strings\r\n")
		}
	}

	s.kvMu.Lock()
	defer s.kvMu.Unlock()

	for _, key := range args {
		delete(s.kv.Strings, key.bulk)
	}
	return []byte("+OK\r\n")
}

func (s *Server) ExistsCmd(args []DataType) []byte {
	fmt.Println("HELLO")
	if len(args) != 1 {
		return []byte("-ERR wrong number of arguments for 'exists' command\r\n")
	}
	key := args[0].bulk
	if key == "" {
		return []byte("-ERR first argument must be a bulk string\r\n")
	}

	s.kvMu.RLock()
	defer s.kvMu.RUnlock()

	s.clearExpiry(key)

	if _, ok := s.kv.Strings[key]; !ok {
		return []byte("#f\r\n")
	} else {
		return []byte("#t\r\n")
	}
}

func (s *Server) IncrCmd(args []DataType) []byte {
	if len(args) != 1 {
		return []byte("-ERR wrong number of arguments for 'incr' command\r\n")
	}
	key := args[0].bulk
	if key == "" {
		return []byte("-ERR first argument must be a bulk string\r\n")
	}

	s.kvMu.Lock()
	defer s.kvMu.Unlock()

	s.clearExpiry(key)

	val := s.kv.Strings[key]
	num, err := strconv.Atoi(val)
	if err != nil {
		return []byte("-ERR value is not a number\r\n")
	}
	num += 1

	s.kv.Strings[key] = strconv.Itoa(num)

	return []byte(":" + strconv.Itoa(num) + "\r\n")
}

func (s *Server) DecrCmd(args []DataType) []byte {
	if len(args) != 1 {
		return []byte("-ERR wrong number of arguments for 'incr' command\r\n")
	}
	key := args[0].bulk
	if key == "" {
		return []byte("-ERR first argument must be a bulk string\r\n")
	}

	s.kvMu.Lock()
	defer s.kvMu.Unlock()

	s.clearExpiry(key)

	val := s.kv.Strings[key]
	num, err := strconv.Atoi(val)
	if err != nil {
		return []byte("-ERR value is not a number\r\n")
	}
	num -= 1

	s.kv.Strings[key] = strconv.Itoa(num)

	return []byte(":" + strconv.Itoa(num) + "\r\n")
}

func (s *Server) ExCmd(args []DataType) []byte {
	if len(args) != 2 {
		return []byte("-ERR wrong number of arguments for 'ex' command\r\n")
	}
	key, secondStr := args[0].bulk, args[1].bulk
	if key == "" || secondStr == "" {
		return []byte("-ERR first two arguments must be bulk strings\r\n")
	}

	seconds, err := strconv.Atoi(secondStr)
	if err != nil {
		return []byte("-ERR second argument must be an integer\r\n")
	}

	s.kvMu.Lock()
	defer s.kvMu.Unlock()

	s.clearExpiry(key)

	if _, ok := s.kv.Strings[key]; !ok {
		return []byte("-ERR key not found\r\n")
	}

	s.kv.Expirations[key] = time.Now().Add(time.Duration(seconds) * time.Second)
	return []byte("+OK\r\n")
}

func (s *Server) TtlCmd(args []DataType) []byte {
	if len(args) != 1 {
		return []byte("-ERR wrong number of arguments for 'TTL' command\r\n")
	}
	key := args[0].bulk
	if key == "" {
		return []byte("-ERR first argument must be bulk strings\r\n")
	}

	s.kvMu.Lock()
	defer s.kvMu.Unlock()

	if _, ok := s.kv.Expirations[key]; !ok {
		return []byte("-ERR key not found\r\n")
	}

	ttl := s.kv.Expirations[key].Sub(time.Now())
	if ttl < 0 {
		delete(s.kv.Strings, key)
		delete(s.kv.Expirations, key)
	}
	return []byte("+" + humanizeDuration(ttl) + "\r\n")
}


func (s *Server) clearExpiry(key string) {
	val, ok := s.kv.Expirations[key]
	if !ok {
		return
	}

	ttl := val.Sub(time.Now())
	if ttl < 0 {
		delete(s.kv.Strings, key)
		delete(s.kv.Expirations, key)
	}
}

func (s *Server) LpushCmd(args []DataType) []byte {
	if len(args) != 2 {
		return []byte("-ERR wrong number of arguments for 'lpush' command\r\n")
	}
	key, val := args[0].bulk, args[1].bulk
	if key == "" || val == "" {
		return []byte("-ERR first two arguments must be bulk strings\r\n")
	}
	s.kvMu.Lock()
	defer s.kvMu.Unlock()

	if _,ok:=s.kv.Lists[key];!ok{
		s.kv.Lists[key]= list.New()
	}
	s.kv.Lists[key].PushFront(val)

	return []byte("+OK\r\n")
}

func (s *Server) RpushCmd(args []DataType) []byte {
	if len(args) != 2 {
		return []byte("-ERR wrong number of arguments for 'rpush' command\r\n")
	}
	key, val := args[0].bulk, args[1].bulk
	if key == "" || val == "" {
		return []byte("-ERR first two arguments must be bulk strings\r\n")
	}
	s.kvMu.Lock()
	defer s.kvMu.Unlock()

	if _,ok:=s.kv.Lists[key];!ok{
		s.kv.Lists[key]= list.New()
	}
	s.kv.Lists[key].PushBack(val)

	return []byte("+OK\r\n")
}

func (s *Server) lpopCmd(args []DataType) []byte {
	if len(args) != 1 {
		return []byte("-ERR wrong number of arguments for 'lpop' command\r\n")
	}
	key := args[0].bulk
	if key == "" {
		return []byte("-ERR first argument must be bulk string\r\n")
	}
	s.kvMu.Lock()
	defer s.kvMu.Unlock()

	if _,ok:=s.kv.Lists[key];!ok{
		return []byte("-ERR key not found\r\n")
	}

	if s.kv.Lists[key].Len()==0{
		return []byte("-ERR list is empty\r\n")
	}

	val := s.kv.Lists[key].Remove(s.kv.Lists[key].Front()).(string)
	return ([]byte("+" + val + "\r\n"))
}

func (s *Server) rpopCmd(args []DataType) []byte {
	if len(args) != 1 {
		return []byte("-ERR wrong number of arguments for 'rpop' command\r\n")
	}
	key := args[0].bulk
	if key == "" {
		return []byte("-ERR first argument must be bulk string\r\n")
	}
	s.kvMu.Lock()
	defer s.kvMu.Unlock()

	if _,ok:=s.kv.Lists[key];!ok{
		return []byte("-ERR key not found\r\n")
	}

	if s.kv.Lists[key].Len()==0{
		return []byte("-ERR list is empty\r\n")
	}

	val := s.kv.Lists[key].Remove(s.kv.Lists[key].Back()).(string)
	return ([]byte("+" + val + "\r\n"))
}
