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
		Atime: time.Unix(stat.Atime, stat.AtimeNsec),
		Ctime: time.Unix(stat.Ctime, stat.CtimeNsec),
	}
}
