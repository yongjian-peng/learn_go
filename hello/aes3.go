package main

import (
	"fmt"

	"example.com/hello/aesCbc"
)

const (
	// CRYPT_KEY_256 = "1~$c31kjtR^@@c2#9&iy"
	CRYPT_KEY_256 = "66313df7e3e7af7ba598caa7d930bede239848b7163440feda68ce3142ebf055"
	CRYPT_KEY_128 = "c31kjtR^@@c2#9&"
)

func main() {
	key := CRYPT_KEY_256
	iv := "123456qo"
	iv = "66313df7e3e7af7ba598caa7d930bede"
	origData := "4rb2auut300790842620672"

	encrData := aesCbc.AesEncrypt([]byte(key), []byte(iv), []byte(origData))

	origData2 := aesCbc.AesDecrypt([]byte(key), []byte(iv), encrData)
	fmt.Println("encrData: ", encrData)
	fmt.Println("origData: ", origData2)

	key = string(encrData) //encode key in bytes to string and keep as secret, put in a vault
	fmt.Printf("key to encrypt/decrypt : %s\n", key)

	key = string(origData2) //encode key in bytes to string and keep as secret, put in a vault
	fmt.Printf("key to encrypt/decrypt : %s\n", key)
}
