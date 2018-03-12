package pgpcrypt

import (
	"bytes"
	"testing"
)

var privKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: Keybase OpenPGP v1.0.0
Comment: https://keybase.io/crypto

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
MwfEz7bN+ypIsO6oX6vthx/pe5Y=
=2HiH
-----END PGP PRIVATE KEY BLOCK-----`

func TestEncrypt(t *testing.T) {
	keyIO := new(bytes.Buffer)
	_, err := keyIO.Write([]byte(privKey))
	if err != nil {
		t.Error(err)
	}

	srcIO := new(bytes.Buffer)
	_, err = srcIO.Write([]byte("test msg"))
	if err != nil {
		t.Error(err)
	}

	dstIO := new(bytes.Buffer)

	err = Encrypt2(keyIO, srcIO, dstIO)
	if err != nil {
		t.Error(err)
	}
}

func TestDecrypt(t *testing.T) {
	keyIO := new(bytes.Buffer)
	_, err := keyIO.Write([]byte(privKey))
	if err != nil {
		t.Error(err)
	}

	srcIO := new(bytes.Buffer)
	_, err = srcIO.Write([]byte("test msg"))
	if err != nil {
		t.Error(err)
	}

	encIO := new(bytes.Buffer)

	err = Encrypt2(keyIO, srcIO, encIO)
	if err != nil {
		t.Error(err)
	}

	keyIO2 := new(bytes.Buffer)
	_, err = keyIO2.Write([]byte(privKey))
	if err != nil {
		t.Error(err)
	}
	decIO := new(bytes.Buffer)
	err = Decrypt2(keyIO2, encIO, decIO, "12345")
	if err != nil {
		t.Error(err)
	}
}
