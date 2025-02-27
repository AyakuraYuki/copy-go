//go:build !(darwin || freebsd || netbsd || plan9 || windows || (wasm && (js || wasip1)))

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
		Atime: time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec)),
		Ctime: time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)),
	}
}
