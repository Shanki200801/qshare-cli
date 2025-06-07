package transfer

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

// SendEncryptedFile reads the file at filePath, encrypts each chunk with key, and sends it over conn.
// If bar is not nil, it updates the progress bar after each chunk.
func SendEncryptedFile(conn io.Writer, filePath string, key []byte, encryptChunk func([]byte, []byte) ([]byte, error), bar *progressbar.ProgressBar) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	buf := make([]byte, 64*1024) // 64KB buffer
	for {
		n, err := file.Read(buf)
		if n > 0 {
			chunk := buf[:n]
			encryptedChunk, err := encryptChunk(key, chunk)
			if err != nil {
				return fmt.Errorf("error encrypting chunk: %w", err)
			}
			if err := binary.Write(conn, binary.BigEndian, uint32(len(encryptedChunk))); err != nil {
				return fmt.Errorf("error writing chunk length: %w", err)
			}
			if _, err := conn.Write(encryptedChunk); err != nil {
				return fmt.Errorf("error sending chunk: %w", err)
			}
			if bar != nil {
				bar.Add(n)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
	}
	binary.Write(conn, binary.BigEndian, uint32(0)) // signal EOF
	return nil
}

// ReceiveAndDecryptFile receives encrypted chunks from conn, decrypts them with key, and writes to outputPath.
// If bar is not nil, it updates the progress bar after each chunk.
func ReceiveAndDecryptFile(conn io.Reader, outputPath string, key []byte, decryptChunk func([]byte, []byte) ([]byte, error), bar *progressbar.ProgressBar) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer out.Close()

	for {
		var chunkLen uint32
		if err := binary.Read(conn, binary.BigEndian, &chunkLen); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading chunk length: %w", err)
		}
		if chunkLen == 0 {
			break // End of file
		}
		encryptedChunk := make([]byte, chunkLen)
		if _, err := io.ReadFull(conn, encryptedChunk); err != nil {
			return fmt.Errorf("error reading chunk: %w", err)
		}
		chunk, err := decryptChunk(key, encryptedChunk)
		if err != nil {
			return fmt.Errorf("error decrypting chunk: %w", err)
		}
		if _, err := out.Write(chunk); err != nil {
			return fmt.Errorf("error writing chunk to file: %w", err)
		}
		if bar != nil {
			bar.Add(len(chunk))
		}
	}
	return nil
}
