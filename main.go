package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path"

	"git.trj.tw/root/go-pgp-tool/pgpcrypt"
	"golang.org/x/crypto/openpgp"

	"git.trj.tw/root/go-pgp-tool/libs"
	"git.trj.tw/root/go-pgp-tool/tools"
)

var (
	flags *libs.Flags
)

func init() {
	flags = new(libs.Flags)
	libs.RegFlag(flags)
	flag.Parse()
}

// args [0] is this
func main() {
	// check flags value
	if !flags.Encrypt && !flags.Decrypt {
		showUsage()
		return
	}
	if flags.Decrypt && flags.Encrypt {
		showUsage()
		return
	}
	if len(flags.KeyFile) == 0 {
		log.Fatal("please input KeyFile path")
	}
	if len(flags.SrcFile) == 0 {
		log.Fatal("please input SrcFile path")
	}
	if len(flags.DstFile) == 0 {
		log.Fatal("please input DstFile path")
	}

	// check file exists
	if !tools.CheckExists(flags.KeyFile, false) {
		log.Fatal("KeyFile not exists")
	}
	if !tools.CheckExists(flags.SrcFile, false) {
		log.Fatal("SrcFile not exists")
	}
	dir := path.Dir(flags.DstFile)
	if !tools.CheckExists(dir, true) {
		log.Fatal("DstFile parent directory not exists")
	}
	if !flags.Override && tools.CheckExists(flags.DstFile, false) {
		log.Fatal("DstFile has Exists if override add flag -y ")
	}

	flags.KeyFile = tools.ParsePath(flags.KeyFile)
	flags.SrcFile = tools.ParsePath(flags.SrcFile)
	flags.DstFile = tools.ParsePath(flags.DstFile)

	if flags.Encrypt {
		// encryptAction()
		encrypt()
	}
	if flags.Decrypt {
		decrypt()
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func showUsage() {
	flag.Usage()
}

func decrypt() {
	keyFile, err := os.Open(flags.KeyFile)
	handleError(err)
	defer keyFile.Close()

	keyList, err := pgpcrypt.ReadKeyFile(keyFile)
	handleError(err)
	if len(keyList) == 0 {
		handleError(errors.New("key not found"))
	}

	key := keyList[0]

	passphraseByte := []byte(flags.Password)
	key.PrivateKey.Decrypt(passphraseByte)
	for _, sub := range key.Subkeys {
		sub.PrivateKey.Decrypt(passphraseByte)
	}

	dstFile, err := os.Create(flags.DstFile)
	handleError(err)
	defer dstFile.Close()

	srcFile, err := os.Open(flags.SrcFile)
	handleError(err)
	defer srcFile.Close()

	md, err := openpgp.ReadMessage(srcFile, keyList, nil, nil)

	_, err = io.Copy(dstFile, md.UnverifiedBody)
	handleError(err)
}

func encrypt() {
	keyFile, err := os.Open(flags.KeyFile)
	handleError(err)
	defer keyFile.Close()

	keyList, err := pgpcrypt.ReadKeyFile(keyFile)
	handleError(err)

	distFile, err := os.Create(flags.DstFile)
	handleError(err)
	defer distFile.Close()

	srcFile, err := os.Open(flags.SrcFile)
	handleError(err)
	defer srcFile.Close()

	fhint := &openpgp.FileHints{}
	fhint.IsBinary = true

	tmpWriter, err := openpgp.Encrypt(distFile, keyList, nil, fhint, nil)
	handleError(err)
	defer tmpWriter.Close()

	_, err = io.Copy(tmpWriter, srcFile)
	handleError(err)

}
