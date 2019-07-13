package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

func DecryptAes128Ecb(data, key []byte) []byte {
	cipher, _ := aes.NewCipher([]byte(key))
	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

func main() {
	f, err := os.Open("data.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := base64.NewDecoder(base64.StdEncoding, f)
	bs, err := ioutil.ReadAll(r)
	key := []byte("YELLOW SUBMARINE")
	decrypted := DecryptAes128Ecb(bs, key)
	fmt.Println(string(decrypted))
}
