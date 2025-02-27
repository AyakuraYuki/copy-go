package copy_go

import (
	"fmt"
	"os"
	"strings"
)

func ExampleCopy() {

	err := Copy(`test/data/example`, `test/data.copy/example`)
	defer func() {
		_ = os.RemoveAll(`test/data.copy`)
	}()
	fmt.Println("Error:", err)

	info, _ := os.Stat(`test/data.copy/example`)
	fmt.Println("IsDir:", info.IsDir())

	// Output:
	// Error: <nil>
	// IsDir: true

}

func ExampleOptions() {

	err := Copy(`test/data/example`, `test/data.copy/example_with_options`, Options{
		Skip: func(src, dst string, info os.FileInfo) (bool, error) {
			return strings.HasSuffix(src, `.git-example`), nil
		},
		OnSymlink: func(string) SymlinkAction {
			return Skip
		},
		PermissionControl: AddPermission(0200),
	})
	defer func() {
		_ = os.RemoveAll(`test/data.copy`)
	}()
	fmt.Println("Error:", err)

	_, err = os.Stat(`test/data.copy/example_with_options/.git-example`)
	fmt.Println("Skipped:", os.IsNotExist(err))

	// Output:
	// Error: <nil>
	// Skipped: true

}
