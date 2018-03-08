package pgpcrypt

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/openpgp"
)

// Encrypt - pgp encrypt func
func Encrypt(key openpgp.EntityList, src io.Reader, dst io.Writer) error {
	fileHint := &openpgp.FileHints{}
	fileHint.IsBinary = true
	fileHint.ModTime = time.Now()

	encWriter, err := openpgp.Encrypt(dst, key, nil, fileHint, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(encWriter, src)
	if err != nil {
		return err
	}

	return nil
}

// EncryptBytes -
func EncryptBytes(key openpgp.EntityList, src io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)

	fileHint := &openpgp.FileHints{}
	fileHint.IsBinary = true

	encWriter, err := openpgp.Encrypt(buf, key, nil, fileHint, nil)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(encWriter, src)
	if err != nil {
		return nil, err
	}

	encBytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return encBytes, nil
}

// Decrypt - pgp decrypt func
func Decrypt(key *openpgp.Entity, keyPassword string, src io.Reader, dst io.Writer) error {
	// decode private key
	passphraseByte := []byte(keyPassword)
	key.PrivateKey.Decrypt(passphraseByte)
	for _, sub := range key.Subkeys {
		sub.PrivateKey.Decrypt(passphraseByte)
	}

	var keyList openpgp.EntityList
	keyList = append(keyList, key)

	md, err := openpgp.ReadMessage(src, keyList, nil, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, md.UnverifiedBody)
	if err != nil {
		return err
	}

	return nil
}
