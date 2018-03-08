package main

import (
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

	// go to decrypt file
	if flags.Decrypt {
		decryptAction()
	}
	if flags.Encrypt {
		// encryptAction()
		encrypt()
	}
}

func decryptAction() {
	// open key file
	keyFile, err := os.Open(flags.KeyFile)
	handleError(err)
	defer keyFile.Close()

	keys, err := pgpcrypt.ReadKeyFile(keyFile)
	handleError(err)
	if len(keys) == 0 {
		log.Fatal("key file not validate")
	}

	srcFile, err := os.Open(flags.SrcFile)
	handleError(err)

	var dstFile *os.File
	if tools.CheckExists(flags.DstFile, false) {
		dstFile, err = os.Open(flags.DstFile)
		handleError(err)
		defer dstFile.Close()
		dstStat, err := dstFile.Stat()
		handleError(err)
		err = dstFile.Truncate(dstStat.Size())
		handleError(err)
	} else {
		dstFile, err = os.Create(flags.DstFile)
		handleError(err)
		defer dstFile.Close()
	}

	key := keys[0]
	err = pgpcrypt.Decrypt(key, flags.Password, srcFile, dstFile)
	handleError(err)
}

func encryptAction() {
	// open key file
	keyFile, err := os.Open(flags.KeyFile)
	handleError(err)
	defer keyFile.Close()

	keys, err := pgpcrypt.ReadKeyFile(keyFile)
	handleError(err)
	if len(keys) == 0 {
		log.Fatal("key file not validate")
	}

	srcFile, err := os.Open(flags.SrcFile)
	handleError(err)

	// encBytes, err := pgpcrypt.EncryptBytes(keys, srcFile)
	// handleError(err)

	// fmt.Println("bytes ::: ", len(encBytes))

	// var dstFile *os.File
	// if tools.CheckExists(flags.DstFile, false) {
	// 	dstFile, err = os.Open(flags.DstFile)
	// 	handleError(err)
	// 	defer dstFile.Close()
	// 	dstStat, err := dstFile.Stat()
	// 	handleError(err)
	// 	err = dstFile.Truncate(dstStat.Size())
	// 	handleError(err)
	// } else {
	dstFile, err := os.Create(flags.DstFile)
	handleError(err)
	defer dstFile.Close()
	// }

	err = pgpcrypt.Encrypt(keys, srcFile, dstFile)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func showUsage() {
	flag.Usage()
}

func encrypt() {
	keyFile, err := os.Open(flags.KeyFile)
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

	distFile, err := os.Create(flags.DstFile)
	if err != nil {
		log.Fatal(err)
	}
	defer distFile.Close()

	// distBuf := new(bytes.Buffer)

	srcFile, err := os.Open(flags.SrcFile)
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
