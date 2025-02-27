//go:build !windows && !plan9

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
		{`~/Music/mora`, fmt.Sprintf("%s/Music/mora", home)},
		{`$HOME/Music/mora`, fmt.Sprintf("%s/Music/mora", home)},
		{`${HOME}/Music/mora`, fmt.Sprintf("%s/Music/mora", home)},
		{`/var/log`, `/var/log`},
		{`/etc/paths.d`, `/etc/paths.d`},
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
		{`~/Music/mora`, fmt.Sprintf("%s/Music/mora", home)},
		{`$HOME/Music/mora`, fmt.Sprintf("%s/Music/mora", home)},
		{`${HOME}/Music/mora`, fmt.Sprintf("%s/Music/mora", home)},
		{`/var/log`, `/var/log`},
		{`/etc/paths.d`, `/etc/paths.d`},
	}

	for _, tt := range tests {
		if get := assureHomeDir(tt.Path); get != tt.Want {
			t.Errorf("assureHomeDir(%q) = %q, want %q", tt.Path, get, tt.Want)
		}
	}
}
