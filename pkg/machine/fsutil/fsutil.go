package fsutil

import (
	"io/fs"
	"os"
)

func FileExists(filepath string) bool {
	return Exists(os.DirFS("/"), filepath)
}

func Exists(osfs fs.FS, fileName string) bool {
	_, err := fs.Stat(osfs, fileName)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return true
}
