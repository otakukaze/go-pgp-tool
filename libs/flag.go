package libs

import (
	"flag"
)

// Flags - flag values struct
type Flags struct {
	Help     bool
	Decrypt  bool
	Encrypt  bool
	SrcFile  string
	DstFile  string
	KeyFile  string
	Override bool
	Password string
}

// RegFlag - Register flag to main
func RegFlag(f *Flags) {
	flag.BoolVar(&f.Help, "h", false, "show usage help")
	flag.BoolVar(&f.Decrypt, "d", false, "decrypt file")
	flag.BoolVar(&f.Encrypt, "e", false, "encrypt file")
	flag.StringVar(&f.SrcFile, "i", "", "input source `file path`")
	flag.StringVar(&f.DstFile, "o", "", "output `file path`")
	flag.StringVar(&f.KeyFile, "k", "", "key `file path`")
	flag.BoolVar(&f.Override, "y", false, "if output file exists override")
	flag.StringVar(&f.Password, "p", "", "private key password")
}
