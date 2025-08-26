package hash

import (
	"github.com/pkg/errors"
	"github.com/speps/go-hashids/v2"
)

// Decode is responsible for decoding a given code.
func Decode(code, alphabet, salt string, minLength int) (uint, error) {
	hashData := hashids.HashIDData{
		Alphabet:  alphabet,
		Salt:      salt,
		MinLength: minLength,
	}

	h, err := hashids.NewWithData(&hashData)
	if err != nil {
		return 0, errors.Wrap(err, "error on creating the hash maker")
	}

	ids, err := h.DecodeWithError(code)
	if err != nil {
		return 0, errors.Wrap(err, "error on decoding hash")
	}

	return uint(ids[0]), nil
}

// Encode is responsible for encoding a given ID.
func Encode(id uint, alphabet, salt string, minLength int) (string, error) {
	hashData := hashids.HashIDData{
		Alphabet:  alphabet,
		Salt:      salt,
		MinLength: minLength,
	}

	h, err := hashids.NewWithData(&hashData)
	if err != nil {
		return "", errors.Wrap(err, "error on creating the hash maker")
	}

	encodedID, err := h.Encode([]int{int(id)})
	if err != nil {
		return "", errors.Wrap(err, "error on encoding hash")
	}

	return encodedID, nil
}
