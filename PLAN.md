# 📍 Development Plan for `qshare`

This file outlines step-by-step milestones and validation points for building `qshare`.

---

## ✅ Phase 0: Setup

- [ ] ✅ Initialize Go module
- [ ] ✅ Setup project folder structure
- [ ] ✅ Add Cobra for CLI

```bash
go mod init github.com/yourusername/qshare
go install github.com/spf13/cobra-cli@latest
cobra-cli init --pkg-name github.com/yourusername/qshare
````

---

## 🛠️ Phase 1: Basic CLI

* [ ] Add `send` and `receive` subcommands
* [ ] Parse `--file` flag in `send` command
* [ ] Accept a code argument in `receive` command

✔️ Validate: Run `./qshare send --file test.txt` and `./qshare receive abc-def-ghi`

---

## 🔄 Phase 2: Code Generator

* [ ] Use 2-3 wordlists + random number to generate one-time codes (e.g. `7-cloud-squid`)
* [ ] Ensure uniqueness for active sessions

✔️ Validate: Printed code format and uniqueness

---

## 🌐 Phase 3: Relay Server (rendezvous only)

* [ ] Create a lightweight TCP or WebSocket relay
* [ ] Map `code -> peer metadata` in memory
* [ ] Allow basic pub/sub for matching codes

✔️ Validate:

* Start relay server
* Sender connects + registers code
* Receiver connects + matches on code

---

## 🔐 Phase 4: Encryption

* [ ] Derive shared key from code (scrypt or PBKDF2)
* [ ] Encrypt file using AES-GCM or NaCl secretbox
* [ ] Add file name and size to encrypted metadata

✔️ Validate: Sender encrypts, receiver decrypts file with correct code

---

## 🚀 Phase 5: Direct P2P File Transfer

* [ ] Use Go `net` package to open TCP socket
* [ ] Use relay server to exchange IP\:port info
* [ ] Establish encrypted P2P session
* [ ] Transfer file over encrypted stream

✔️ Validate: Transfer file over LAN with direct connection

---

## 📦 Phase 6: Chunked Transfer + Resume

* [ ] Break files into chunks (e.g., 1MB)
* [ ] Send chunk checksums for verification
* [ ] Support resume after disconnection

✔️ Validate: Mid-transfer resume works correctly

---

## 🌍 Phase 7 (Optional): NAT Traversal

* [ ] Add UPnP / NAT hole punching
* [ ] Fallback to relay file transfer

---

## 🧪 Phase 8: Polish

* [ ] Add TUI progress bar
* [ ] Improve logging and error handling
* [ ] Add timeout and retry logic

---

## 🧼 Phase 9: Docs & Packaging

* [ ] Polish README.md
* [ ] Create usage examples
* [ ] Optional: Publish binary via GitHub Releases
