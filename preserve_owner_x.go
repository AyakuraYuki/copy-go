//go:build windows || plan9

package copy_go

import "io/fs"

func preserveOwner(src, dst string, info fs.FileInfo) (err error) {
	return nil // do nothing
}
