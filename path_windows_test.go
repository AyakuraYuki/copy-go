//go:build windows

package copy_go

import (
	"fmt"
	"testing"
)

func Test_expandHomeDir(t *testing.T) {
	home, err := homeDir()
	if err != nil {
		t.Fatalf("cannot get home dir: %v", err)
	}

	tests := []struct {
		Path string
		Want string
	}{
		{`%userprofile%\Music\mora`, fmt.Sprintf(`%s\Music\mora`, home)},
		{`C:\Users\Temp`, `C:\Users\Temp`},
	}

	for _, tt := range tests {
		if get := expandHomeDir(tt.Path); get != tt.Want {
			t.Errorf("expandHomeDir(%q) = %q, want %q", tt.Path, get, tt.Want)
		}
	}
}

func Test_assureHomeDir(t *testing.T) {
	home, err := homeDir()
	if err != nil {
		t.Fatalf("cannot get home dir: %v", err)
	}

	tests := []struct {
		Path string
		Want string
	}{
		{`%userprofile%\Music\mora`, fmt.Sprintf(`%s\Music\mora`, home)},
		{`C:\Users\Temp`, `C:\Users\Temp`},
	}

	for _, tt := range tests {
		if get := assureHomeDir(tt.Path); get != tt.Want {
			t.Errorf("assureHomeDir(%q) = %q, want %q", tt.Path, get, tt.Want)
		}
	}
}
