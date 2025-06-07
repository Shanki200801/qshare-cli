# Go Progress Bar Primer (schollz/progressbar)

## Why Use a Progress Bar?
- Shows users how much of a file has been sent/received.
- Improves UX for large/slow transfers.

---

## Popular Library: schollz/progressbar
- URL: https://github.com/schollz/progressbar
- Easy to use, supports dynamic updates, works in any terminal.
- Install: `go get github.com/schollz/progressbar/v3`

---

## Basic Usage
```go
import "github.com/schollz/progressbar/v3"

bar := progressbar.Default(totalBytes) // totalBytes = file size
for each chunk {
    // ... send or receive chunk ...
    bar.Add(len(chunk))
}
```
- If you don't know the total size, use `progressbar.Default(-1)` for an indeterminate bar.

---

## Example: File Send with Progress Bar
```go
fileInfo, _ := os.Stat(filePath)
bar := progressbar.Default(fileInfo.Size())
for {
    n, err := file.Read(buf)
    if n > 0 {
        // ... send chunk ...
        bar.Add(n)
    }
    // ...
}
```

## Example: File Receive with Progress Bar
- If you know the file size (e.g., sent as metadata), use it for the bar.
- Otherwise, use an indeterminate bar.

```go
bar := progressbar.Default(-1)
for {
    // ... receive chunk ...
    bar.Add(len(chunk))
}
```

---

## Python Equivalent
```python
from tqdm import tqdm
for chunk in tqdm(chunks, total=total_chunks):
    # send/receive chunk
```

## Java Equivalent
- No standard library, but you can print progress manually or use third-party libraries.

---

## Tips
- Always update the bar after each chunk.
- For best UX, print the bar to stderr (so stdout is clean for piping).
- For receive, consider sending the file size as metadata for a determinate bar.

---

## Summary
- schollz/progressbar is the go-to for Go CLI progress bars.
- Easy to integrate into chunked file transfer.
- Use total size if known, otherwise use indeterminate mode. 