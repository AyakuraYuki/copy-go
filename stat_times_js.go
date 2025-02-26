//go:build js

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
		Atime: time.Unix(int64(stat.Atime), int64(stat.AtimeNsec)),
		Ctime: time.Unix(int64(stat.Ctime), int64(stat.CtimeNsec)),
	}
}
