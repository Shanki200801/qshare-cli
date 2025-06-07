# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

qshare is a secure, peer-to-peer encrypted file sharing CLI tool built in Go. It enables users to transfer files directly between machines using human-readable codes (e.g., "5-apple-7-tiger") with end-to-end encryption.

## Development Commands

### Build
```bash
# Build the main CLI tool
go build -o qshare .

# Build the relay server
go build -o relay-server ./relay-server/relay-server.go
```

### Run
```bash
# Start the relay server (required for file transfers)
./relay-server

# Send a file
./qshare send --file path/to/file.txt

# Receive a file
./qshare receive 7-tiger-cloud
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/crypto
```

### Dependencies
```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify

# Tidy dependencies
go mod tidy
```

## Architecture

### Component Overview

The application follows a modular architecture with clear separation of concerns:

1. **Main CLI (`main.go`)**: Cobra-based CLI interface with `send` and `receive` commands
2. **Relay Server (`relay-server/`)**: TCP server that facilitates peer discovery and connection
3. **Internal Packages**: Core functionality organized by domain

### Internal Package Structure

- **`codegen`**: Generates memorable sharing codes using pattern `{digit}-{fruit}-{digit}-{animal}`
- **`crypto`**: AES-GCM encryption with SHA-256 key derivation from sharing codes
- **`transfer`**: Chunked file transfer (64KB chunks) with encryption, supports directory zipping
- **`validate`**: Input validation for file paths and operations
- **`relay`**: Connection management and security rate limiting

### Data Flow

1. **Sender Flow**:
   - Generate unique code via `codegen.GenerateCode()`
   - Connect to relay server on `localhost:4000`
   - Register with code and wait for receiver
   - Encrypt file in chunks using derived key
   - Stream encrypted chunks to receiver

2. **Receiver Flow**:
   - Connect to relay server with provided code
   - Derive decryption key from code
   - Receive encrypted chunks
   - Decrypt and save to disk

### Security Features

- **End-to-End Encryption**: AES-GCM with 256-bit keys
- **Rate Limiting**: Prevents brute force attacks (5 attempts/minute per IP/code)
- **Failed Handshake Protection**: Blocks codes after 3 failed attempts
- **One-Time Codes**: Each transfer uses a unique code
- **Optional Extra Key**: Additional entropy via `-ekey` flag

### Key Design Decisions

- **Chunked Transfer**: 64KB chunks for efficient memory usage and progress tracking
- **Directory Support**: Automatic zipping/unzipping with metadata preservation
- **Retry Support**: Optional reconnection within 2 minutes for interrupted transfers
- **Binary Protocol**: BigEndian encoding for chunk length prefixes
- **Channel-Based Relay**: Go channels for bidirectional peer communication