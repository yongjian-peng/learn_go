package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func main() {

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}
	fmt.Printf("key to bytes : %s\n", bytes)
	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	// fmt.Printf("key to encrypt/decrypt : %s\n", key)

	key = "62953a6997071c99a4d9c09385a79966490f4cfac8e2c253d8bfc67b98f89fda"
	fmt.Printf("key to encrypt/decrypt : %s\n", key)
	encrypted := encrypt("payout_token=Y3pTRVVUR2hTN1NwSHNrM3dLMUlpM3JXNi9PS3BPRndlSzQ3c1dmSkJHTkxLOU9ycTFzSHMvL2xqNVBGZHBHKw&payout_secret_key=Needsix@0007 &checkout=CURL&client_ip=192.168.101.1&source=Encode-CurlAPI&source_url=https://needsixgaming.com¶m1=param1¶m2=param2¶m3=param¬ify_url=notify_url&success_url=success_url&error_url=error_url", key)
	fmt.Printf("encrypted : %s\n", encrypted)

	decrypted := decrypt(encrypted, key)
	fmt.Printf("decrypted : %s\n", decrypted)
}

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {
	fmt.Println("keyString: ", keyString)
	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)

	fmt.Println("keyddd: ", key)

	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
