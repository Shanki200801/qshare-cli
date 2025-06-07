# Go Primer: CLI UX Improvements for File Transfer

## Why Good UX?
- Makes the CLI tool intuitive and pleasant for users.
- Reduces friction for common tasks (e.g., sending directories, handling errors).

---

## Example: Auto-Zipping Directories
```go
info, err := os.Stat(filePath)
if info.IsDir() {
    zipPath, err := transfer.ZipDir(filePath)
    filePath = zipPath
    defer os.Remove(zipPath)
}
```
**What does this do?**
- Checks if the path is a directory. If so, zips it to a temp file before sending.
- Uses Go's `os.Stat`, `archive/zip`, and `filepath.WalkDir`.
- `defer os.Remove(zipPath)` ensures the temp file is deleted after use.

---

## Example: Progress Bar with schollz/progressbar
```go
import "github.com/schollz/progressbar/v3"

bar := progressbar.Default(fileInfo.Size())
for {
    n, err := file.Read(buf)
    if n > 0 {
        // ... send chunk ...
        bar.Add(n)
    }
}
```
**What does this do?**
- Shows a progress bar for file transfer, updating after each chunk.
- If file size is unknown, use `progressbar.Default(-1)` for indeterminate mode.

---

## Example: Retryable Transfers (Sender)
```go
sendCmd.Flags().BoolVarP(&allowRetry, "allowRetry", "r", false, "Allow sender to reconnect within 2 minutes if disconnected during transfer")
if allowRetry {
    fmt.Fprintf(conn, "%s:sender:retry\n", code)
} else {
    fmt.Fprintf(conn, "%s:sender\n", code)
}
```
**What does this do?**
- Adds a CLI flag to enable retryable transfers.
- If enabled, the sender handshake includes `:retry` so the relay server allows reconnects.

---

## Example: User Feedback on Errors
```go
if codeFromHandshake != "" {
    allowed, triesLeft, blocked, blockMsg := relay.CheckAndRecordFailedHandshake(codeFromHandshake)
    if !allowed {
        if blocked {
            conn.Write([]byte(blockMsg + "\n"))
        } else {
            msg := fmt.Sprintf("Invalid code or key. You have %d tries remaining before this code is blocked.\n", triesLeft)
            conn.Write([]byte(msg))
        }
    }
}
```
**What does this do?**
- Tells the user how many tries remain before a code is blocked, or if the code is temporarily blocked.

---

## Go Keywords & Libraries
- `os.Stat`, `defer`, `archive/zip`, `filepath.WalkDir`, `github.com/schollz/progressbar/v3`, `flag`, `fmt.Fprintf`, `if`, `for`, `import`

---

## Key UX Features
- **Auto-Zipping Directories:** If a directory is provided, it is zipped to a temp file before sending. (Uses Go's archive/zip and filepath.WalkDir.)
- **Progress Bar:** Uses schollz/progressbar to show transfer progress (determinate if file size known, indeterminate otherwise).
- **Retryable Transfers:** Sender can enable a retry window; receiver always supports retry if sender does.
- **User Feedback:** Clear error messages, including how many tries remain before a code is blocked, and disconnect notifications.

---

## Go Concepts Used
- **Flag Parsing:** With spf13/cobra for CLI arguments.
- **os.Stat:** To check if a path is a file or directory.
- **Defer/Cleanup:** Ensures temp files are deleted after use.
- **Progress Bar:** Updated after each chunk sent/received.

---

## Keywords
- CLI, UX, progress bar, auto-zip, error feedback, retry, temp file, flag parsing, defer. 