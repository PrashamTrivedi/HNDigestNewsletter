package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

//RegistrationData holds data of registration
type RegistrationData struct {
	Email string
	Types []string
}

func init() {
	fmt.Println("This should be printed first")
}

func main() {

	encryptedData := encrypt([]byte(`{"email":"prash2488@gmail.com","types":["show","ask","jobs"]}`), "m@yth3f0rc3w1t4y0u")
	fmt.Println(string(encryptedData))

	decryptedData := decrypt(encryptedData, "m@yth3f0rc3w1t4y0u")
	fmt.Println(string(decryptedData))
	var registrationData RegistrationData
	err := json.Unmarshal(decryptedData, &registrationData)
	if err != nil {
		panic(err)
	}
	fmt.Println(registrationData)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}
