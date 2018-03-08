package tools

import (
	"os"
	"path"
	"runtime"
	"strings"
)

// ParsePath - parse file path to absPath
func ParsePath(wd, dst string) string {
	if []rune(dst)[0] == '~' {
		home := UserHomeDir()
		if len(home) > 0 {
			dst = strings.Replace(dst, "~", home, -1)
		}
	}

	if path.IsAbs(dst) {
		dst = path.Clean(dst)
		return dst
	}

	str := path.Join(wd, dst)
	str = path.Clean(str)
	return str
}

// UserHomeDir - get user home directory
func UserHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

// CheckExists - check file exists
func CheckExists(filePath string, allowDir bool) bool {
	stat, err := os.Stat(filePath)
	if err != nil && !os.IsExist(err) {
		return false
	}
	if !allowDir && stat.IsDir() {
		return false
	}
	return true
}
