package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/otakukaze/go-pgp-tool/pgpcrypt"

	"github.com/otakukaze/go-pgp-tool/libs"
	"github.com/otakukaze/go-pgp-tool/tools"
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
		// encrypt()
		keyIO, err := os.Open(flags.KeyFile)
		handleError(err)

		srcIO, err := os.Open(flags.SrcFile)
		handleError(err)

		dstIO, err := os.Create(flags.DstFile)
		handleError(err)

		err = pgpcrypt.Encrypt2(keyIO, srcIO, dstIO)

		// err := pgpcrypt.Encrypt(flags.KeyFile, flags.SrcFile, flags.DstFile)
		handleError(err)
	}
	if flags.Decrypt {
		// decrypt()
		keyIO, err := os.Open(flags.KeyFile)
		handleError(err)

		srcIO, err := os.Open(flags.SrcFile)
		handleError(err)

		dstIO, err := os.Create(flags.DstFile)
		handleError(err)

		err = pgpcrypt.Decrypt2(keyIO, srcIO, dstIO, flags.Password)
		// err := pgpcrypt.Decrypt(flags.KeyFile, flags.SrcFile, flags.DstFile, flags.Password)
		handleError(err)
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
