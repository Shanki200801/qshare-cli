# Go Encryption Primer: AES-GCM & Key Derivation

## Symmetric Encryption in Go
- Go's standard library provides crypto/aes for AES encryption.
- AES-GCM (Galois/Counter Mode) is secure and provides both confidentiality and integrity.

## Key Derivation
- You should never use a passphrase/code directly as an AES key.
- Use a key derivation function (KDF) like PBKDF2 or scrypt to turn a code into a 32-byte key.
- Go provides crypto/sha256 for hashing and golang.org/x/crypto/pbkdf2 for PBKDF2.

---

## Example: Deriving a Key from a Code
```go
import (
    "crypto/sha256"
)

func deriveKey(code string) []byte {
    hash := sha256.Sum256([]byte(code))
    return hash[:]
}
```
- This is a simple hash-based KDF (good enough for a CLI tool).
- For more security, use PBKDF2 with a salt and many iterations.

---

## Example: Encrypting with AES-GCM
```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
)

func encrypt(key, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}
```
- Nonce is random and prepended to the ciphertext.
- Use the same key and nonce to decrypt.

---

## Example: Decrypting with AES-GCM
```go
func decrypt(key, ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}
```

---

## Python Equivalent
```python
from Crypto.Cipher import AES
from hashlib import sha256
key = sha256(code.encode()).digest()
cipher = AES.new(key, AES.MODE_GCM)
nonce = cipher.nonce
ciphertext, tag = cipher.encrypt_and_digest(data)
```

## Java Equivalent
```java
SecretKeySpec key = new SecretKeySpec(MessageDigest.getInstance("SHA-256").digest(code.getBytes()), "AES");
Cipher cipher = Cipher.getInstance("AES/GCM/NoPadding");
cipher.init(Cipher.ENCRYPT_MODE, key);
byte[] ciphertext = cipher.doFinal(data);
```

---

## Summary
- Derive a 32-byte key from the code (hash or PBKDF2).
- Use AES-GCM for encryption/decryption.
- Prepend the nonce to the ciphertext for decryption.
- This is secure and easy to implement in Go. 