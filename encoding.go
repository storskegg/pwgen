package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/awnumar/memguard"
)

const dictionaryX = "0123456789abcdefghijklmno-qrstuvwxyzABCDEF_HIJKLMNOPQRSTUVWXYZ"
const dictHex = "0123456789abcdef"

// Encode is a wrapper to encode slice b given the dictionaryX constant
func Encode(pwLen int, hex bool) string {
	var enc *Encoding
	var err error
	if hex {
		enc, err = NewEncoding(dictHex)
	} else {
		enc, err = NewEncoding(dictionaryX)
	}
	if err != nil {
		fmt.Println(err)
		memguard.SafeExit(2)
	}

	return enc.Encode(pwLen)
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

// Encode function receives a byte slice and encodes it to a string using the alphabet provided
func (e *Encoding) Encode(maxLen int) string {
	source, err := memguard.NewImmutableRandom(maxLen)
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

	result := res.String()

	return result[0:maxLen]
}
