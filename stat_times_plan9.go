package copy_go

import (
	"os"
)

// todo: check plan9 in the future

func getTimeSpec(info os.FileInfo) timespec {
	return timespec{
		Mtime: info.ModTime(),
		Atime: info.ModTime(),
		Ctime: info.ModTime(),
	}
}
