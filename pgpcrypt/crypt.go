package pgpcrypt

import (
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/openpgp"
)

// Decrypt2 - decrypt
func Decrypt2(keyReader io.Reader, srcReader io.Reader, dstReader io.Writer, password string) (err error) {
	keyList, err := ReadKeyFile(keyReader)
	if err != nil {
		return
	}

	if len(keyList) == 0 {
		return errors.New("key not found")
	}

	key := keyList[0]
	passphraseByte := []byte(password)
	key.PrivateKey.Decrypt(passphraseByte)
	for _, sub := range key.Subkeys {
		sub.PrivateKey.Decrypt(passphraseByte)
	}

	md, err := openpgp.ReadMessage(srcReader, keyList, nil, nil)

	if md == nil {
		return errors.New("key password error")
	}

	_, err = io.Copy(dstReader, md.UnverifiedBody)
	if err != nil {
		return err
	}

	return
}

// Decrypt - decrypt file with pgp
func Decrypt(KeyFile, SrcFile, DstFile, Password string) error {
	keyFile, err := os.Open(KeyFile)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	keyList, err := ReadKeyFile(keyFile)
	if err != nil {
		return err
	}
	if len(keyList) == 0 {
		return errors.New("key not found")
	}

	key := keyList[0]

	passphraseByte := []byte(Password)
	key.PrivateKey.Decrypt(passphraseByte)
	for _, sub := range key.Subkeys {
		sub.PrivateKey.Decrypt(passphraseByte)
	}

	dstFile, err := os.Create(DstFile)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := os.Open(SrcFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	md, err := openpgp.ReadMessage(srcFile, keyList, nil, nil)

	if md == nil {
		return errors.New("key password error")
	}

	_, err = io.Copy(dstFile, md.UnverifiedBody)
	if err != nil {
		return err
	}

	return nil
}

// Encrypt2 - encrypt
func Encrypt2(keyReader io.Reader, srcReader io.Reader, dstWriter io.Writer) (err error) {
	keyList, err := ReadKeyFile(keyReader)
	if err != nil {
		return
	}

	fhint := &openpgp.FileHints{}
	fhint.IsBinary = true

	tmpWriter, err := openpgp.Encrypt(dstWriter, keyList, nil, fhint, nil)
	if err != nil {
		return
	}
	defer tmpWriter.Close()

	_, err = io.Copy(tmpWriter, srcReader)
	if err != nil {
		return
	}

	return
}

// Encrypt - encrypt file with pgp
func Encrypt(KeyFile, SrcFile, DstFile string) error {
	keyFile, err := os.Open(KeyFile)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	keyList, err := ReadKeyFile(keyFile)
	if err != nil {
		return err
	}

	distFile, err := os.Create(DstFile)
	if err != nil {
		return err
	}
	defer distFile.Close()

	srcFile, err := os.Open(SrcFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	fhint := &openpgp.FileHints{}
	fhint.IsBinary = true

	tmpWriter, err := openpgp.Encrypt(distFile, keyList, nil, fhint, nil)
	if err != nil {
		return err
	}
	defer tmpWriter.Close()

	_, err = io.Copy(tmpWriter, srcFile)
	if err != nil {
		return err
	}

	return nil
}
