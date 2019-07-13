package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	src := []byte(`49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d`)
	h := make([]byte, hex.DecodedLen(len(src)))
	hex.Decode(h, src)
	fmt.Printf("Source string: %s\n", h)
	fmt.Println(base64.StdEncoding.EncodeToString(h), "encoded")
	fmt.Println(`SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t`, "expected")
}
