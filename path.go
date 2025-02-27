package copy_go

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func assureHomeDir(path string) string {
	path = expandHomeDir(path)
	ret, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return ret
}

func expandHomeDir(path string) string {
	home, err := homeDir()
	if err != nil || home == "" {
		return path
	}

	old := ""
	switch runtime.GOOS {
	case "windows":
		if strings.HasPrefix(strings.ToLower(path), "%userprofile%") {
			old = "%userprofile%"
		}

	case "plan9":
		if strings.HasPrefix(path, "$home/") {
			old = "$home"
		}

	default:
		if strings.HasPrefix(path, "~/") {
			old = "~"
		}
		if strings.HasPrefix(path, "$HOME/") {
			old = "$HOME"
		}
		if strings.HasPrefix(path, "${HOME}/") {
			old = "${HOME}"
		}
	}

	if old != "" {
		path = strings.Replace(path, old, home, 1)
	}
	return path
}

func homeDir() (string, error) {
	cur, err := user.Current()
	if err == nil && cur.HomeDir != "" {
		return cur.HomeDir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get determine home directory via environment: %w", err)
	}
	return home, nil
}
