package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/awnumar/memguard"
	"math"
)

// TODO: Will need some pad for arbitrary lengths in the future

const dictionary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const magic = 1.34375

// Encode is a wrapper to encode slice b given the dictionary constant
func Encode() string {
	enc, err := NewEncoding(dictionary)
	if err != nil {
		fmt.Println(err)
		memguard.SafeExit(2)
	}

	return enc.Encode()
}

// The following has been taken from https://github.com/eknkc/basex/blob/master/basex.go and minimally adapted

// Encoding is a custom base encoding defined by an alphabet.
// It should bre created using NewEncoding function
type Encoding struct {
	base        int
	alphabet    []rune
	alphabetMap map[rune]int
}

// NewEncoding returns a custom base encoder defined by the alphabet string.
// The alphabet should contain non-repeating characters.
// Ordering is important.
// Example alphabets:
//   - base2: 01
//   - base16: 0123456789abcdef
//   - base32: 0123456789ABCDEFGHJKMNPQRSTVWXYZ
//   - base62: 0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
func NewEncoding(alphabet string) (*Encoding, error) {
	runes := []rune(alphabet)
	runeMap := make(map[rune]int)

	for i := 0; i < len(runes); i++ {
		if _, ok := runeMap[runes[i]]; ok {
			return nil, errors.New("ambiguous alphabet")
		}

		runeMap[runes[i]] = i
	}

	return &Encoding{
		base:        len(runes),
		alphabet:    runes,
		alphabetMap: runeMap,
	}, nil
}

func bufSize() int {
	return int(math.Floor(float64(pwLen) / magic))
}

// Encode function receives a byte slice and encodes it to a string using the alphabet provided
func (e *Encoding) Encode() string {
	source, err := memguard.NewImmutableRandom(bufSize())
	if err != nil {
		fmt.Println(err)
		memguard.SafeExit(2)
	}

	defer source.Destroy()

	if source.Size() == 0 {
		return ""
	}

	digits := []int{0}

	for i := 0; i < source.Size(); i++ {
		carry := int(source.Buffer()[i])

		for j := 0; j < len(digits); j++ {
			carry += digits[j] << 8
			digits[j] = carry % e.base
			carry = carry / e.base
		}

		for carry > 0 {
			digits = append(digits, carry%e.base)
			carry = carry / e.base
		}
	}

	var res bytes.Buffer

	for k := 0; source.Buffer()[k] == 0 && k < source.Size()-1; k++ {
		res.WriteRune(e.alphabet[0])
	}

	for q := len(digits) - 1; q >= 0; q-- {
		res.WriteRune(e.alphabet[digits[q]])
	}

	return res.String()
}
