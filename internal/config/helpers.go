package config

import (
	"log/slog"
	"os"
	"os/user"
	"path/filepath"

	"github.com/liangzhaoliang95/lxz/internal/slogs"
)

// IsBoolSet checks if a bool ptr is set.
func IsBoolSet(b *bool) bool {
	return b != nil && *b
}

// isEnvSet checks if env var is set.
func isEnvSet(env string) bool {
	return os.Getenv(env) != ""
}

// UserTmpDir returns the temp dir with the current user name.
func UserTmpDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(os.TempDir(), u.Username, AppName)

	return dir, nil
}

// MustLXZUser establishes current user identity or fail.
func MustLXZUser() string {
	usr, err := user.Current()
	if err != nil {
		envUsr := os.Getenv("USER")
		if envUsr != "" {
			return envUsr
		}
		envUsr = os.Getenv("LOGNAME")
		if envUsr != "" {
			return envUsr
		}
		slog.Error("Die on retrieving user info", slogs.Error, err)
		os.Exit(1)
	}
	return usr.Username
}
