# Go Concurrency Primer

## 1. Goroutines
A goroutine is a lightweight thread managed by Go.
You start one by putting `go` before a function call:

```go
go fmt.Println("Hello from a goroutine!")
```
- This runs the function in the background.
- Your program keeps running; the goroutine runs "concurrently."

---

## 2. Channels
Channels let goroutines communicate safely.
Think of them as pipes for passing data.

```go
ch := make(chan string) // create a channel for strings

go func() {
    ch <- "hello" // send "hello" into the channel
}()

msg := <-ch // receive from the channel
fmt.Println(msg) // prints "hello"
```
- `<-` is the "arrow" for sending/receiving.
- Channels block: if you try to receive and nothing's there, you wait.

---

## 3. Select
`select` lets you wait on multiple channel operations at once.

```go
select {
case msg := <-ch1:
    fmt.Println("got", msg, "from ch1")
case msg := <-ch2:
    fmt.Println("got", msg, "from ch2")
}
```
- Whichever channel is ready first, that case runs.

---

## 4. Closing Channels
When you're done sending, close the channel:

```go
close(ch)
```
- Receivers can check if a channel is closed.

---

## 5. Example: Simple Producer/Consumer

```go
ch := make(chan int)

go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}()

for n := range ch {
    fmt.Println(n)
}
```
- The goroutine sends numbers 0–4, then closes the channel.
- The main goroutine receives and prints them.

---

## Why does this matter for your project?
- You'll use goroutines to handle sender and receiver "sessions."
- You'll use channels (or similar) to simulate a relay server: sender waits for receiver, receiver connects using the code. 