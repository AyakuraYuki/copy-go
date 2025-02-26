//go:build windows || plan9 || js

package copy_go

func preserveLtimes(src, dst string) (err error) {
	return nil
}
