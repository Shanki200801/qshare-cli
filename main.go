package main

import (
	"fmt"
	"net"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/shanki200801/qshare/internal/codegen"
	"github.com/shanki200801/qshare/internal/crypto"
	"github.com/shanki200801/qshare/internal/transfer"
	"github.com/shanki200801/qshare/validate"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "qshare",
		Short: "qshare is a p2p file sharing CLI tool",
	}

	var filePath string
	var ekey string
	var outputPath string

	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send a file",
		Run: func(comd *cobra.Command, args []string) {
			// Validate the file exists and is not a directory
			if err := validate.ValidateFile(filePath); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if ekey != "" {
				fmt.Println("Using encryption key:", ekey)
			}
			// Connect to the relay server
			conn, err := net.Dial("tcp", "localhost:4000")
			if err != nil {
				fmt.Println("Error connecting to relay server:", err)
				os.Exit(1)
			}
			defer conn.Close()
			// Generate and print the one-time code
			code := codegen.GenerateCode()
			fmt.Println("Your code is:", code)
			// Handshake: identify as sender
			fmt.Fprintf(conn, "%s:sender\n", code)
			// Derive encryption key from code and ekey
			key := crypto.DeriveKey(code, ekey)
			// Create a progress bar for file transfer
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Println("Error getting file info:", err)
				os.Exit(1)
			}
			bar := progressbar.Default(fileInfo.Size())
			// Send the file in encrypted chunks with progress bar
			if err := transfer.SendEncryptedFile(conn, filePath, key, crypto.Encrypt, bar); err != nil {
				fmt.Println("Error sending file:", err)
				os.Exit(1)
			}
			fmt.Println("File sent successfully")
		},
	}
	sendCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the file to send")
	sendCmd.Flags().StringVar(&ekey, "ekey", "", "Extra encryption key (must match on receive)")
	sendCmd.MarkFlagRequired("file")

	var receiveCmd = &cobra.Command{
		Use:   "receive []",
		Short: "Receive a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			code := args[0]
			if ekey != "" {
				fmt.Println("Using encryption key:", ekey)
			}
			// Derive decryption key from code and ekey
			key := crypto.DeriveKey(code, ekey)
			// Connect to the relay server
			conn, err := net.Dial("tcp", "localhost:4000")
			if err != nil {
				fmt.Println("Error connecting to relay server:", err)
				os.Exit(1)
			}
			defer conn.Close()
			// Handshake: identify as receiver
			fmt.Fprintf(conn, "%s:receiver\n", code)
			// Use an indeterminate progress bar (file size unknown)
			bar := progressbar.Default(-1)
			// Receive and decrypt the file in chunks with progress bar
			if err := transfer.ReceiveAndDecryptFile(conn, outputPath, key, crypto.Decrypt, bar); err != nil {
				fmt.Printf("Error receiving file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("File received and decrypted successfully! Saved as: %s\n", outputPath)
		},
	}
	receiveCmd.Flags().StringVarP(&outputPath, "output", "o", "Received_file", "Output file path")
	receiveCmd.Flags().StringVar(&ekey, "ekey", "", "Extra encryption key (must match sender)")

	rootCmd.AddCommand(sendCmd, receiveCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
