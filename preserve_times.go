package copy_go

import "os"

func preserveTimes(dst string, info os.FileInfo) error {
	spec := getTimeSpec(info)
	if err := os.Chtimes(dst, spec.Atime, spec.Mtime); err != nil {
		return err
	}
	return nil
}
