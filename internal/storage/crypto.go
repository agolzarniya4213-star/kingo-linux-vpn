package storage

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
    "os"
    "path/filepath"
)

type CryptoLayer struct {
    key []byte
}

func NewCryptoLayer(keyPath string) (*CryptoLayer, error) {
    // Ensure directory exists
    os.MkdirAll(filepath.Dir(keyPath), 0700)

    key, err := os.ReadFile(keyPath)
    if err != nil {
        if os.IsNotExist(err) {
            // Generate new 32-byte key for AES-256
            key = make([]byte, 32)
            if _, err := rand.Read(key); err != nil {
                return nil, err
            }
            if err := os.WriteFile(keyPath, key, 0600); err != nil {
                return nil, err
            }
        } else {
            return nil, err
        }
    }

    return &CryptoLayer{key: key}, nil
}

func (c *CryptoLayer) Encrypt(plaintext string) string {
    block, err := aes.NewCipher(c.key)
    if err != nil {
        return plaintext // Fallback
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return plaintext
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return plaintext
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext)
}

func (c *CryptoLayer) Decrypt(encodedText string) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encodedText)
    if err != nil {
        // Fallback for old plaintext data (Migration)
        return encodedText, nil
    }

    block, err := aes.NewCipher(c.key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        // Fallback for old plaintext data
        return encodedText, nil
    }

    return string(plaintext), nil
}
