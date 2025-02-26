//go:build aix || illumos || js || netbsd || plan9 || solaris || windows

package copy_go

import (
	"os"
)

// windows does not support named pipes
func pcopy(dst string, info os.FileInfo) error {
	return nil
}
