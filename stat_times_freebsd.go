//go:build freebsd

package copy_go

import (
	"os"
	"syscall"
	"time"
)

func getTimeSpec(info os.FileInfo) timespec {
	stat := info.Sys().(*syscall.Stat_t)
	return timespec{
		Mtime: info.ModTime(),
		Atime: time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec),
		Ctime: time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec),
	}
}
