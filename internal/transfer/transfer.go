package transfer

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"archive/zip"
	"io/fs"
	"path/filepath"

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

// ZipDir zips the contents of srcDir into a temp zip file and returns the path to the zip file.
// Preserves file permissions and timestamps. Caller is responsible for deleting the temp file after use.
func ZipDir(srcDir string) (string, error) {
	tmpFile, err := os.CreateTemp("", "qshare-*.zip")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	zipPath := tmpFile.Name()
	err = ZipDirTo(srcDir, zipPath)
	tmpFile.Close()
	if err != nil {
		os.Remove(zipPath)
		return "", err
	}
	return zipPath, nil
}

// ZipDirTo zips the contents of srcDir into destZip, preserving structure, permissions, and timestamps.
func ZipDirTo(srcDir, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return fmt.Errorf("error creating zip file: %w", err)
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // skip directories, only add files
		}
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = relPath
		hdr.Method = zip.Deflate
		fWriter, err := zipWriter.CreateHeader(hdr)
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(fWriter, f)
		return err
	})
}

// Unzip extracts a zip archive at zipPath into destDir, preserving structure, permissions, and timestamps.
func Unzip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("error opening zip file: %w", err)
	}
	defer r.Close()
	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
		if f.Modified.IsZero() == false {
			os.Chtimes(fpath, f.Modified, f.Modified)
		}
	}
	return nil
}
