# Go Function Definition Syntax

## 1. Basic Function
```go
func FunctionName(params) returnType {
    // function body
}
```
- `func` — keyword to define a function.
- `FunctionName` — name of the function (capitalized = exported/public, lowercase = private).
- `params` — comma-separated list of parameters (name type).
- `returnType` — type(s) of value(s) returned (can be omitted if nothing is returned).

**Example:**
```go
func Add(a int, b int) int {
    return a + b
}
```

---

## 2. Method (Function with a Receiver)
```go
func (r *Relay) CreateRoom(code string) chan []byte {
    // method body
}
```
- `(r *Relay)` — receiver: this method is attached to the Relay type.  
  - `r` is the variable name for the receiver (like self/this in other languages).
  - `*Relay` means it operates on a pointer to Relay (so it can modify the struct).
- `CreateRoom` — method name.
- `code string` — parameter (name type).
- `chan []byte` — return type (a channel of byte slices).

**Why use a receiver?**
- Methods with receivers can access and modify the struct's fields (like mu and rooms in Relay).

---

## 3. Multiple Return Values
Go functions can return multiple values, often (value, error) or (value, ok).

```go
func (r *Relay) JoinRoom(code string) (chan []byte, bool) {
    // returns a channel and a boolean
}
```

---

## 4. Exported vs Unexported
- Capitalized function/method names (e.g., CreateRoom) are exported (public).
- Lowercase names (e.g., newRelay) are unexported (private to the package). 