package copy_go

import (
	"io/fs"
	"os"
)

// tmpDirectoryWritablePermission makes the destination directory writable,
// so that stuff can be copied recursively even if any original directory is NOT writable.
const tmpDirectoryWritablePermission = os.FileMode(0755)

type PermissionControlFunc func(srcinfo fs.FileInfo, dst string) (chmodfunc func(*error), err error)

var (
	AddPermission = func(perm os.FileMode) PermissionControlFunc {
		return func(srcinfo fs.FileInfo, dst string) (func(*error), error) {
			orig := srcinfo.Mode()
			if srcinfo.IsDir() {
				if err := os.MkdirAll(dst, tmpDirectoryWritablePermission); err != nil {
					return func(*error) {}, err
				}
			}
			return func(err *error) {
				chmod(dst, orig|perm, err)
			}, nil
		}
	}
	PreservePermission                       = AddPermission(0)
	DoNothing          PermissionControlFunc = func(srcinfo fs.FileInfo, dst string) (func(*error), error) {
		if srcinfo.IsDir() {
			if err := os.MkdirAll(dst, srcinfo.Mode()); err != nil {
				return func(*error) {}, err
			}
		}
		return func(e *error) {}, nil
	}
)

// chmod ANYHOW changes file mode,
// with assigning error raised during Chmod
// BUT respecting the error already reported.
func chmod(dir string, mode os.FileMode, reported *error) {
	if err := os.Chmod(dir, mode); *reported == nil {
		*reported = err
	}
}
