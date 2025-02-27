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
		Atime: time.Unix(0, int64(stat.Atime)),
		Ctime: time.Unix(0, int64(stat.Ctime)),
	}
}
