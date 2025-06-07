# Qshare

A secure, peer-to-peer encrypted file sharing CLI tool inspired by [magic-wormhole](https://github.com/magic-wormhole/magic-wormhole), built in Go.

## ✨ Features

- 🔐 End-to-end encrypted file transfer
- ⚡ Peer-to-peer direct connection (relay fallback supported)
- 🔑 Easy-to-share one-time code (e.g. `5-sky-train`)
- 📦 Chunked file transfer with integrity checks
- 🧪 Simple, terminal-based CLI

## 🧰 Tech Stack

- Language: Go
- CLI: [Cobra](https://github.com/spf13/cobra)
- Networking: net, optional libp2p or custom TCP/QUIC
- Encryption: NaCl/libsodium (via `golang.org/x/crypto/nacl/secretbox` or AES-GCM)

## 🚀 Getting Started

### Install

```bash
git clone https://github.com/yourusername/qshare.git
cd qshare
go build -o qshare .
````

### Usage

#### Send a file

```bash
./qshare send --file path/to/file.txt
# Outputs: Your code is: 7-tiger-cloud
```

#### Receive a file

```bash
./qshare receive 7-tiger-cloud
# Downloads and saves the file securely
```

## 📦 Architecture Overview

1. **Sender** starts a session and generates a code
2. **Receiver** uses the code to rendezvous via a lightweight relay server
3. Once matched, both peers establish a secure connection
4. The file is encrypted and sent directly over the wire

## 🔒 Security

* Code-derived key exchange (PBKDF2 or ECDH)
* File encryption with NaCl/AES
* Encrypted metadata + chunks
* One-time use codes for security

## 🛣️ Roadmap

* [ ] Basic CLI with send/receive commands
* [ ] In-memory relay server for code matching
* [ ] Encrypted file transfer (single file)
* [ ] Chunked transfer + resume support
* [ ] NAT traversal (UPnP/STUN)
* [ ] Optional relay file transfer fallback