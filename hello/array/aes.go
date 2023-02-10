package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	plaintext := []byte("12345678912345678912345678900000")
	key, err := base64.StdEncoding.DecodeString("cidRgzwfcgztwae/mccalIeedOAmA/CbU3HEqWz1Ejk=")
	fmt.Println("key: ", key)
	if err != nil {
		panic(err)
	}
	iv, err := hex.DecodeString("162578ddce177a4a7cb2f7c738fa052d")
	if err != nil {
		panic(err)
	}

	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(ciphertext, plaintext)

	fmt.Printf("EncryptedText %v\n", string(ciphertext))
	fmt.Printf("EncryptedText as hex %v\n", hex.EncodeToString(ciphertext))
	fmt.Printf("EncryptedText as base 64 %v\n", base64.StdEncoding.EncodeToString(ciphertext))
}

func pkcs7Pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
