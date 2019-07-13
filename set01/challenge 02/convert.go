package convert

import "fmt"

// FixedXOR returns two equal-length byte arrays XORed with each other
func FixedXOR(b1, b2 []byte) ([]byte, error) {
	if len(b1) != len(b2) {
		return nil, fmt.Errorf("FixedXOR needs equal length byte arrays")
	}
	res := make([]byte, len(b1))
	for i, v1 := range b1 {
		res[i] = v1 ^ b2[i]
	}
	return res, nil
}
