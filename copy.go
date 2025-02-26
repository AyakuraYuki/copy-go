package copy_go

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sync/semaphore"
)

// Copy copies src to dst, no matter if src is a file or a directory
func Copy(src, dst string, opts ...Options) (err error) {
	opt := assureOptions(src, dst, opts...)

	if opt.NumOfWorkers > 1 {
		opt.intent.ctx = context.Background()
		opt.intent.sem = semaphore.NewWeighted(opt.NumOfWorkers)
	}

	var info fs.FileInfo
	if opt.FS != nil {
		info, err = fs.Stat(opt.FS, src)
	} else {
		info, err = os.Lstat(src)
	}
	if err != nil {
		return onError(src, dst, err, opt)
	}

	return switchboard(src, dst, info, opt)
}

// switchboard switches proper copy functions regarding file type, etc...
// If there would be anything else here, add a case to this switchboard.
func switchboard(src, dst string, info os.FileInfo, opt Options) (err error) {
	if info.Mode()&os.ModeDevice != 0 && !opt.Specials {
		return onError(src, dst, err, opt)
	}

	if opt.RenameDestination != nil {
		if dst, err = opt.RenameDestination(src, dst); err != nil {
			return onError(src, dst, err, opt)
		}
	}

	switch {
	case info.Mode()&os.ModeSymlink != 0:
		err = onSymlink(src, dst, opt)
	case info.Mode()&os.ModeNamedPipe != 0:
		err = pcopy(dst, info)
	case info.IsDir():
		err = dcopy(src, dst, info, opt)
	default:
		err = fcopy(src, dst, info, opt)
	}

	return onError(src, dst, err, opt)
}

// fcopy is for just a file,
// with considering existence of parent directory and file permission.
func fcopy(src, dst string, info os.FileInfo, opt Options) (err error) {
	var readCloser io.ReadCloser
	if opt.FS != nil {
		readCloser, err = opt.FS.Open(src)
	} else {
		readCloser, err = os.Open(src)
	}
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer fclose(readCloser, &err)

	if err = os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fclose(f, &err)

	chmodfunc, err := opt.PermissionControl(info, dst)
	if err != nil {
		return err
	}
	chmodfunc(&err)

	var reader = bufio.NewReader(readCloser)
	var writer = bufio.NewWriter(f)
	if opt.CopyBufferSize > 0 {
		reader = bufio.NewReaderSize(reader, opt.CopyBufferSize)
		writer = bufio.NewWriterSize(writer, opt.CopyBufferSize)
	}

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	if opt.Sync {
		err = f.Sync()
	}

	if opt.PreserveOwner {
		if err := preserveOwner(src, dst, info); err != nil {
			return err
		}
	}

	if opt.PreserveTimes {
		if err := preserveTimes(dst, info); err != nil {
			return err
		}
	}

	return err
}

func dcopy(src, dst string, info os.FileInfo, opt Options) (err error) {

}

func onSymlink(src, dst string, opt Options) error {
	switch opt.OnSymlink(src) {
	case Deep:
		orig, err := os.Readlink(src)
		if err != nil {
			return err
		}
		if !filepath.IsAbs(orig) {
			orig = filepath.Join(filepath.Dir(src), orig) // orig is a relative link, need to concat src dir
		}
		info, err := os.Lstat(orig)
		if err != nil {
			return err
		}
		return copyNextOrSkip(orig, dst, info, opt)

	case Shallow:
		if err := lcopy(src, dst); err != nil {
			return err
		}
		if opt.PreserveTimes {
			return preserveLtimes(src, dst)
		}
		return nil

	case Skip:
		fallthrough

	default:
		return nil // do nothing, act not supported
	}
}

// copyNextOrSkip decides if this src should be copied or not.
// because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
func copyNextOrSkip(src, dst string, info os.FileInfo, opt Options) error {
	if opt.Skip != nil {
		skip, err := opt.Skip(src, dst, info)
		if err != nil {
			return err
		}
		if skip {
			return nil
		}
	}
	return switchboard(src, dst, info, opt)
}

// lcopy is for a symlink, with just creating a new symlink by replicating src symlink
func lcopy(src, dst string) error {
	orig, err := os.Readlink(src)
	// ** might be controlled by Options in the future **
	if err != nil {
		if os.IsNotExist(err) {
			return os.Symlink(src, dst) // copy symlink even if not existing
		}
		return err
	}

	// ** might be controlled by SymlinkExistsAction **
	if _, err = os.Lstat(dst); err == nil {
		if err = os.Remove(dst); err != nil {
			return err
		}
	}

	return os.Symlink(orig, dst)
}

// fclose ANYHOW closes file,
// with assigning error raised during Close,
// BUT respecting the error already reported.
func fclose(f io.Closer, reported *error) {
	if err := f.Close(); *reported == nil {
		*reported = err
	}
}

// onError lets caller handle errors occurred when copying
func onError(src, dst string, err error, opt Options) error {
	if opt.OnError == nil {
		return err
	}
	return opt.OnError(src, dst, err)
}
