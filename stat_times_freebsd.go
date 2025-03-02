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
		Atime: time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec)),
		Ctime: time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec)),
	}
}
