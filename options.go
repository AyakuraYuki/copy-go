package copy_go

import (
	"context"
	"io/fs"
	"os"

	"golang.org/x/sync/semaphore"
)

type Options struct {
	// OnSymlink can specify what to do on symlink
	OnSymlink func(src string) SymlinkAction

	// OnDirExists can specify what to do when there's a directory already existing in destination
	OnDirExists func(src, dst string) DirExistsAction

	// OnError lets caller decide whether to continue on particular copy error
	OnError func(src, dst string, err error) error

	// Skip can specify which files should be skipped
	Skip func(src, dst string, srcinfo os.FileInfo) (bool, error)

	// RenameDestination can specify the destination file or dir name if needed to rename
	RenameDestination func(src, dst string) (string, error)

	// Specials includes special files to be copied (default: false)
	Specials bool

	// AddPermission to every entity
	// DO NOT MORE THAN 0777
	// @OBSOLETE
	// Use `PermissionControl = AddPermission(perm)` instead
	AddPermission os.FileMode

	// PermissionControl can preserve or even add permission to every entity.
	// for example:
	//
	//		opt.PermissionControl = AddPermission(0222)
	//
	// see `permission.go` for more detail
	PermissionControl PermissionControlFunc

	// Sync file after copy.
	// Useful in case when file must be on the disk
	// (in case crash happens, for example),
	// at the expense of some performance penalty
	Sync bool

	// PreserveOwner preserve the uid and the gid of all entries
	PreserveOwner bool

	// PreserveTimes preserve the atime and the mtime of the entries.
	// On linux we can preserve only up to 1 millisecond accuracy.
	PreserveTimes bool

	// The byte size of the buffer to use for copying files.
	// Leave it to zero to use the default buffer size.
	CopyBufferSize int

	// If given, copy.Copy refers to this fs.FS instead of the OS filesystem.
	// e.g., You can use embed.FS to copy files from embedded filesystem.
	FS fs.FS

	// NumOfWorkers represents the number of workers used for
	// concurrent copying contents of directories.
	// If 0 or 1, it does not use goroutine for copying directories.
	// Please refer to https://pkg.go.dev/golang.org/x/sync/semaphore for more details.
	NumOfWorkers int64

	// PreferConcurrent is a function to determine whether
	// to use goroutine for copying contents of directories.
	// If PreferConcurrent is nil, which is default, it does concurrent
	// copying for all directories.
	// If NumOfWorkers is 0 or 1, this function will be ignored.
	PreferConcurrent func(src, dst string) (bool, error)

	// internal use only
	intent intent
}

type intent struct {
	src string
	dst string
	ctx context.Context
	sem *semaphore.Weighted
}

type SymlinkAction int

const (
	Deep    SymlinkAction = iota // Deep creates hard-copy of contents
	Shallow                      // Shallow creates new symlink to the dst of symlink
	Skip                         // Skip does nothing with symlink
)

type DirExistsAction int

const (
	Merge       DirExistsAction = iota // Merge preserves or overwrites existing files under the dir (default behavior)
	Replace                            // Replace deletes all contents under the dir and copy src files
	Untouchable                        // Untouchable does nothing for the dir, and leaves it as it is
)

// getDefaultOptions provides default options
func getDefaultOptions(src, dst string) Options {
	return Options{
		OnSymlink: func(string) SymlinkAction {
			return Shallow // default: do shallow copy
		},
		OnDirExists:       nil,                // default: Merge
		OnError:           nil,                // default: accept error
		Skip:              nil,                // default: do NOT skip
		RenameDestination: nil,                // default: no rename
		Specials:          false,              // default: do NOT copy special files
		AddPermission:     0,                  // default: add nothing
		PermissionControl: PreservePermission, // default: just preserve permission
		Sync:              false,              // default: do NOT sync
		PreserveOwner:     false,              // default: do NOT preserve owner
		PreserveTimes:     false,              // default: do NOT preserve the modification time
		CopyBufferSize:    0,                  // default: use default buffer size
		FS:                nil,                // default: do not specify file system
		NumOfWorkers:      0,                  // default: copy in sequential
		PreferConcurrent:  nil,                // default: no concurrent
		intent: intent{
			src: src,
			dst: dst,
			ctx: nil,
			sem: nil,
		},
	}
}

// assureOptions struct, should be called only once
// all optional values MUST NOT BE nil/zero after assured
func assureOptions(src, dst string, opts ...Options) Options {
	defaults := getDefaultOptions(src, dst)
	if len(opts) == 0 {
		return defaults
	}
	if opts[0].OnSymlink == nil {
		opts[0].OnSymlink = defaults.OnSymlink
	}
	if opts[0].Skip == nil {
		opts[0].Skip = defaults.Skip
	}
	if opts[0].AddPermission > 0 {
		opts[0].PermissionControl = AddPermission(opts[0].AddPermission)
	} else if opts[0].PermissionControl == nil {
		opts[0].PermissionControl = PreservePermission
	}
	opts[0].intent.src = defaults.intent.src
	opts[0].intent.dst = defaults.intent.dst
	return opts[0]
}

func shouldCopyDirectoryConcurrent(src, dst string, opt Options) (bool, error) {
	if opt.NumOfWorkers <= 1 {
		return false, nil
	}
	if opt.PreferConcurrent == nil {
		return true, nil
	}
	return opt.PreferConcurrent(src, dst)
}
