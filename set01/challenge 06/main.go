package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func HammingDistanceString(s1, s2 string) int {
	var d int
	b2 := []rune(s2)
	for i, r1 := range []rune(s1) {
		d += bits.OnesCount32(uint32(r1 ^ b2[i]))
	}
	return d
}

func HammingDistanceBytes(b1, b2 []byte) int {
	var d int
	for i, v1 := range b1 {
		d += bits.OnesCount8(uint8(v1 ^ b2[i]))
	}
	return d
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

type keySizeCand struct {
	// key size candidate
	KeySize  int
	Distance float32
}

func GuessKeySize(data []byte, keepBest int) []keySizeCand {
	maxKeySize := 40
	keySizeCands := make([]keySizeCand, 0, maxKeySize-1)
	for keySize := 2; keySize <= maxKeySize; keySize++ {
		blocks := len(data) / keySize / 2
		if blocks == 0 {
			break
		}
		// sum up normalized hamming distance over all consecutive pairs of blocks
		var sum float32
		for block := 0; block < blocks; block++ {
			i1 := 2 * block * keySize
			i2 := i1 + keySize
			hd := HammingDistanceBytes(data[i1:i2], data[i2:i2+keySize])
			sum += float32(hd) / float32(keySize)
		}
		keySizeCands = append(keySizeCands, keySizeCand{keySize, sum / float32(blocks)})
	}
	sort.Slice(keySizeCands, func(i, j int) bool {
		return keySizeCands[i].Distance < keySizeCands[j].Distance
	})
	if len(keySizeCands) < keepBest {
		return keySizeCands
	}
	return keySizeCands[:keepBest]
}

func GuessRepeatingXorKey(data []byte) {
	for _, cand := range GuessKeySize(data, 5) {
		key := make([]rune, cand.KeySize)
		for ki := 0; ki < cand.KeySize; ki++ {
			testBuf := make([]byte, 0, len(data)/cand.KeySize)
			for i := ki; i < len(data); i += cand.KeySize {
				testBuf = append(testBuf, data[i])
			}
			ch, _, _ := GuessXoredValue(testBuf)
			key[ki] = ch
		}
		fmt.Printf("cand %+v key %s\n", cand, string(key))
	}
}

func main() {
	f, err := os.Open("data.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := base64.NewDecoder(base64.StdEncoding, f)
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	log.Println("read", len(bs), "bytes from file")
	//s1 := "this is a test"
	//s2 := "wokka wokka!!!"
	//fmt.Println("Hamming dist (expected 37):", HammingDistance(s1, s2))
	GuessRepeatingXorKey(bs)
}
