package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"golang.org/x/crypto/openpgp"
)

// args [0] is this
func main() {
	if len(os.Args) < 3 {
		log.Fatal("encrypt keyPath filePath")
	}
	keyPath := os.Args[1]
	filePath := os.Args[2]

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	keyPath = path.Join(wd, keyPath)
	filePath = path.Join(wd, filePath)

	if _, err := os.Stat(keyPath); err != nil && !os.IsExist(err) {
		fmt.Println(err)
		log.Fatal("key file not exists")
	}

	if _, err := os.Stat(filePath); err != nil && !os.IsExist(err) {
		log.Fatal("key file not exists")
	}

	encrypt(keyPath, filePath)
}

func encrypt(pubKey, srcPath string) {
	keyFile, err := os.Open(pubKey)
	if err != nil {
		log.Fatal(err)
	}
	defer keyFile.Close()

	var keyList openpgp.EntityList
	keys, err := openpgp.ReadArmoredKeyRing(keyFile)
	// keys, err := openpgp.ReadKeyRing(keyFile)
	if err != nil {
		log.Fatal(err)
	}

	keyList = append(keyList, keys...)

	distFile, err := os.Create("./dist.pgp")
	if err != nil {
		log.Fatal(err)
	}
	defer distFile.Close()

	// distBuf := new(bytes.Buffer)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	fhint := &openpgp.FileHints{}
	fhint.IsBinary = true

	tmpWriter, err := openpgp.Encrypt(distFile, keyList, nil, fhint, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tmpWriter.Close()

	_, err = io.Copy(tmpWriter, srcFile)
	if err != nil {
		log.Fatal(err)
	}

}
