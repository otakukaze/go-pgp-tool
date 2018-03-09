package pgpcrypt

import (
	"io"

	"golang.org/x/crypto/openpgp"
)

// ReadKeyFile - read key from file
func ReadKeyFile(r io.Reader) (openpgp.EntityList, error) {
	keys, err := openpgp.ReadArmoredKeyRing(r)
	if err != nil || len(keys) == 0 {
		keys, err = openpgp.ReadKeyRing(r)
		if err != nil {
			return nil, err
		}
	}

	return keys, nil
}

// CombineKeys - combine key
func CombineKeys(keys ...openpgp.EntityList) openpgp.EntityList {
	var keyList openpgp.EntityList
	for _, key := range keys {
		keyList = append(keyList, key...)
	}

	return keyList
}
