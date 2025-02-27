package copy_go

import (
	"time"
)

type timespec struct {
	Mtime time.Time // modify time
	Atime time.Time // access time
	Ctime time.Time // change time
}
