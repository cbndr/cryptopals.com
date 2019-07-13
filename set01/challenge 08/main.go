package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
)

func decodeHexString(src string) []byte {
	h := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(h, []byte(src))
	if err != nil {
		panic(err)
	}
	return h
}

func countIdenticalChunks(data []byte) int {
	const chunkSize = 16
	m := map[string]int{}
	max := len(data) - chunkSize
	for i := 0; i < max; i += chunkSize {
		m[string(data[i:i+chunkSize])]++
	}
	var best int
	for _, score := range m {
		if score > best {
			best = score
		}
	}
	return best
}

// https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#ECB

func main() {
	f, err := os.Open("data.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var lineNo int
	var bestLineNo, bestScore int
	// find the line with most repeating 16-byte chunks in it
	for sc.Scan() {
		lineNo++
		// count identical 16 byte chunks
		if score := countIdenticalChunks(decodeHexString(sc.Text())); score > bestScore {
			bestScore = score
			bestLineNo = lineNo
		}
	}
	fmt.Println("Result: line no", bestLineNo, "with score", bestScore)
}
