package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func decodeBytes(src string) []byte {
	h := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(h, []byte(src))
	if err != nil {
		panic(err)
	}
	return h
}

func xorByVal(bs []byte, b byte) []byte {
	res := make([]byte, len(bs))
	for i, b1 := range bs {
		res[i] = b1 ^ b
	}
	return res
}

// englishLetterFrequency returns the english letter frequency of a string
func englishLetterFrequency(input string) float32 {
	// http://en.algoritmy.net/article/40379/Letter-frequency-English
	var englishFreq = []float32{
		0.08167, 0.01492, 0.02782, 0.04253, 0.12702, 0.02228, 0.02015, // A-G
		0.06094, 0.06966, 0.00153, 0.00772, 0.04025, 0.02406, 0.06749, // H-N
		0.07507, 0.01929, 0.00095, 0.05987, 0.06327, 0.09056, 0.02758, // O-U
		0.00978, 0.02360, 0.00150, 0.01974, 0.00074,                   // V-Z
	}
	var observed [26]int
	var ignored, invalid int

	for _, ch := range input {
		if ch >= 65 && ch <= 90 {
			observed[ch-65]++ // uppercase A-Z
		} else if ch >= 97 && ch <= 122 {
			observed[ch-97]++ // lowercase a-z
		} else if strings.ContainsRune("\"'`´ .,:;", ch) {
			ignored++
		} else {
			invalid++
		}
	}
	res := float32(1) - float32(invalid)/float32(len(input))
	charCount := float32(len(input) - ignored)
	for i := 0; i < 26; i++ {
		res += float32(observed[i])/charCount - englishFreq[i]
	}
	return res
}

func GuessXoredValue(input []byte) (rune, float32, string) {
	var bestFrequ float32
	var bestRune rune
	for r := rune(1); r < 128; r++ {
		test := string(xorByVal(input, byte(r)))
		if frequ := englishLetterFrequency(test); frequ > bestFrequ {
			bestFrequ = frequ
			bestRune = r
		}
	}
	return bestRune, bestFrequ, string(xorByVal(input, byte(bestRune)))
}

func main() {
	f, err := os.Open("sentences.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var lineNo = 1
	var maxFrequ float32
	var vs []interface{}
	for sc.Scan() {
		ch, frequ, s := GuessXoredValue(decodeBytes(sc.Text()))
		/*if frequ > 0.8 {
			fmt.Println("Candidate", lineNo, "XOR char", string(ch), "hit rate", frequ, "decoded:", s)
		}*/
		if frequ > maxFrequ {
			vs = []interface{}{"Guessing line", lineNo, "XOR char", string(ch), "hit rate", frequ, "decoded:", s}
			maxFrequ = frequ
		}
		lineNo++
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	fmt.Println(vs...)
}
