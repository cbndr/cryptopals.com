package convert

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func decodeBytes(src string) []byte {
	h := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(h, []byte(src))
	if err != nil {
		panic(err)
	}
	return h
}

func TestFixedXOR(t *testing.T) {
	h1 := decodeBytes(`1c0111001f010100061a024b53535009181c`)
	h2 := decodeBytes(`686974207468652062756c6c277320657965`)
	result, err := FixedXOR(h1, h2)
	if err != nil {
		panic(err)
	}
	expected := decodeBytes(`746865206b696420646f6e277420706c6179`)
	if !bytes.Equal(result, expected) {
		t.Errorf("bytes are not equal:\n'%x'\n'%x'", string(expected), string(result))
	}
}
