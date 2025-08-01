# NekoDB

A lightweight in-memory key-value store server implemented in Go that communicates using the Redis Serialization Protocol (RESP). It supports strings, lists, and sets along with expiration, concurrency safety, and basic command execution.

## ✨ Features

- RESP protocol support
- String operations: `SET`, `GET`, `DEL`, `EXISTS`, `INCR`, `DECR`
- Expiry commands: `EX`, `TTL`
- List commands: `LPUSH`, `RPUSH`, `LPOP`, `RPOP`
- Set commands: `SADD`, `SREM`, `SMEMBERS`
- Basic commands: `PING`, `ECHO`
- Concurrent-safe in-memory storage using `sync.RWMutex`
- Periodic cleanup of expired keys

## 🛠️ Getting Started

### Prerequisites

- Go 1.18+
- `make` (optional, for Makefile support)
#### 📦 Dependencies
* charmbracelet/log – Elegant structured logger for CLI apps

Install dependencies:
```bash
go mod tidy
```
### Installation

Clone the repository:

```bash
git clone https://github.com/vikram-parashar/nekodb/
cd nekodb
```

#### Build & Run:
(Using the Makefile)
```bash
make run
```
(or manually)
```bash
go build -o kv-server
./kv-server
```
Sever will start at localhost:8080

# 🔌 Using redis-cli
You can test the server using the official Redis command-line client (redis-cli), which understands RESP.

## Run redis-cli against your server:
```bash
redis-cli -p 8080
```
### Example Commands:
```bash
127.0.0.1:8080> PING
"PONG"

127.0.0.1:8080> SET name vikram
"OK"

127.0.0.1:8080> GET name
"vikram"

127.0.0.1:8080> EX name 5
"OK"

127.0.0.1:8080> TTL name
"4s"
```
# 🧩 Project Structure
```bash
.
├── main.go           # Server entry point
├── server.go         # TCP server logic
├── cmds.go           # RESP command implementations
├── parser.go         # RESP protocol parser
├── helperFunc.go     # Utilities (TTL formatting, cleanup)
├── go.mod/go.sum     # Module dependencies
├── Makefile          # Build/run targets
```

Made with ❤️ in Go
