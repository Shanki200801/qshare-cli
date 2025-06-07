# Go Primer: Secure Code Generation for File Transfer

## Why Code Entropy Matters
- Higher entropy (randomness) makes codes harder to guess or brute-force.
- Prevents attackers from easily discovering valid transfer codes.

---

## Example: New Code Format
```go
func GenerateCode() string {
    num1 := rand.Intn(9000) + 1000 // 1000-9999
    num2 := rand.Intn(9000) + 1000
    w1 := words[rand.Intn(len(words))]
    w2 := words[rand.Intn(len(words))]
    return fmt.Sprintf("%d-%s-%d-%s", num1, w1, num2, w2)
}
```
**What does this do?**
- Picks two random 4-digit numbers and two random words from a list.
- Combines them into a code like `1234-apple-5678-zebra`.
- Uses Go's `rand.Intn` for randomness and `fmt.Sprintf` for formatting.

---

## Security Impact
- Makes brute-force/code-guessing attacks extremely impractical, especially when combined with rate limiting and lockout.
- Each code has much higher entropy than a short or single-word code.

---

## Go Keywords & Libraries
- `rand.Intn`, `fmt.Sprintf`, `slice`, `string formatting`, `import math/rand`, `import fmt`

---

## Go Concepts Used
- **rand.Intn:** For generating random numbers.
- **Random Indexing:** For picking words from a slice.
- **fmt.Sprintf:** For formatting the final code string.

---

## Keywords
- Entropy, randomness, brute-force prevention, code generation, rand, slice, string formatting. 