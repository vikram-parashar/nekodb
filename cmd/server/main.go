package main

import (
	"nekodb/internal/server"
)

func main() {
	server := server.NewServer(":8080")
	server.Start()
}
