# NekoDB

[![Go Version](https://img.shields.io/badge/go-1.18+-brightgreen)](https://golang.org/doc/go1.18)
[![License](https://img.shields.io/github/license/vikram-parashar/nekodb)](./LICENSE)
[![Build Status](https://github.com/vikram-parashar/nekodb/actions/workflows/go.yml/badge.svg)](https://github.com/vikram-parashar/nekodb/actions)

A lightweight, in-memory key-value database, compatible with most Redis clients. NekoDB is designed for simplicity, fast prototyping, and local development where spinning up a full Redis instance is overkill.

---

## Features

- 🐾 **Redis Protocol Compatible:** Use your favorite Redis client libraries out of the box.
- 💾 **In-memory Storage:** Fast, ephemeral, and ideal for testing or non-persistent scenarios.
- 🔒 **Concurrent Access:** Handles multiple simultaneous client connections.
- ⏰ **Key Expiry:** Supports expiring keys automatically.
- 🧹 **Automatic Cleanup:** Periodically removes expired keys to free memory.
- 🗃️ **Data Types:** Supports strings, lists, sets, and more (see below).
- ⚡ **Zero Dependencies:** Single binary, no external services required.

---

## Getting Started

### Installation

```sh
go install github.com/vikram-parashar/nekodb@latest
```

Or clone and build from source:

```sh
git clone https://github.com/vikram-parashar/nekodb.git
cd nekodb
go build -o nekodb .
```

### Usage

Start the server (default port is 6379):

```sh
./nekodb
```

Or specify a port:

```sh
./nekodb -port 6380
```

Connect with your favorite Redis client:

```python
# Python example using redis-py
import redis
r = redis.Redis(host='localhost', port=6379)
r.set('foo', 'bar')
print(r.get('foo'))  # b'bar'
```

---

## Supported Commands & Data Types

| Command    | Description                                     |
|------------|-------------------------------------------------|
| SET/GET    | Set and get string values                       |
| DEL        | Delete a key                                    |
| EXPIRE     | Set key expiration                              |
| LPUSH/RPUSH| List operations                                 |
| SADD/SMEMBERS | Set operations                               |
| KEYS       | List all keys (with pattern)                    |

> **Note:** Some advanced Redis commands are not supported. See [docs/COMMANDS.md](docs/COMMANDS.md) for a full list.

---

## Configuration

- `-port` — Specify the listening port (default: 6379)
- `-cleanup-interval` — Set key expiration cleanup interval (default: 60s)

---

## Contributing

Contributions are welcome! Please submit issues or pull requests via GitHub.

1. Fork the repo
2. Create a feature branch
3. Push your changes
4. Open a pull request

---

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) for details.

---

## Acknowledgements

- Inspired by [Redis](https://redis.io/)
- Built with [Go](https://golang.org/)

---

## Contact & Support

- For issues, feature requests, or questions, please use [GitHub Issues](https://github.com/vikram-parashar/nekodb/issues).
