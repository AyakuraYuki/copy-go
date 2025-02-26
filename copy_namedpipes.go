//go:build !aix && !illumos && !js && !netbsd && !plan9 && !solaris && !windows

package copy_go

import (
	"os"
	"path/filepath"
	"syscall"
)

func pcopy(dst string, info os.FileInfo) error {
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}
	return syscall.Mkfifo(dst, uint32(info.Mode()))
}
