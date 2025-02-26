//go:build !windows && !plan9 && !js

package copy_go

import "golang.org/x/sys/unix"

func preserveLtimes(src, dst string) (err error) {
	info := new(unix.Stat_t)
	if err = unix.Lstat(src, info); err != nil {
		return err
	}
	return unix.Lutimes(dst, []unix.Timeval{
		unix.NsecToTimeval(info.Atim.Nano()),
		unix.NsecToTimeval(info.Mtim.Nano()),
	})
}
