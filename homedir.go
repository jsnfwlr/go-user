package user

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// HomeDirCache controls caching of the home directory. Caching is enabled
// by default.
var HomeDirCache = true

var (
	dirCacheValue string
	dirCacheLock  sync.RWMutex
)

// Reset clears the cache, forcing the next call to Dir to re-detect
// the home directory. This generally never has to be called, but can be
// useful in tests if you're modifying the home directory via the HOME
// env var or something.
func ResetHomeDir() {
	dirCacheLock.Lock()
	defer dirCacheLock.Unlock()
	dirCacheValue = ""
}

// HomeDir returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func HomeDir() (value string, fault error) {
	if HomeDirCache {
		dirCacheLock.RLock()
		cached := dirCacheValue
		dirCacheLock.RUnlock()
		if cached != "" {
			return cached, nil
		}
	}

	dirCacheLock.Lock()
	defer dirCacheLock.Unlock()

	result, err := dir()
	if err != nil {
		return "", err
	}
	dirCacheValue = result
	return result, nil
}

// ExpandPath expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
func ExpandPath(path string) (value string, fault error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != os.PathSeparator {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := HomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

// dirMethods is a list of functions that can be used to detect the home
// directory. The first method to return a non-empty string is used.
// The actual functions are added to the slice in the OS-specific files.
var dirMethods = []func() string{}

func dir() (string, error) {
	for _, m := range dirMethods {
		if dir := m(); dir != "" {
			return dir, nil
		}
	}

	return "", errors.New("could not determine home directory")
}

// homeEnv is the environment variable that should the home directory for
// most systems.
var homeEnv = "HOME"

// OS agnostic, base-level method to get the home directory. This is the
// first method to be tried. The value for homeEnv is determined in the
// OS-specific files.
func homeEnvVar() string {
	return os.Getenv(homeEnv)
}

// Mac OS/Darwin
func dirDSCL() string {
	var buff bytes.Buffer
	cmd := exec.Command("sh", "-c", `dscl -q . -read /Users/"$(whoami)" NFSHomeDirectory | sed 's/^[^ ]*: //'`)
	cmd.Stdout = &buff
	if err := cmd.Run(); err == nil {
		result := strings.TrimSpace(buff.String())
		if result != "" {
			return result
		}
	}
	return ""
}

// Unix like systems
func cdPwd() string {
	// If all else fails, try the shell
	var buff bytes.Buffer
	cmd := exec.Command("sh", "-c", "cd && pwd")
	cmd.Stdout = &buff
	if err := cmd.Run(); err != nil {
		return ""
	}

	result := strings.TrimSpace(buff.String())
	if result == "" {
		return ""
	}

	return result
}

// plan9, linux
func dirGetEntPasswd() string {
	var buff bytes.Buffer

	cmd := exec.Command("getent", "passwd", strconv.Itoa(os.Getuid()))
	cmd.Stdout = &buff
	if err := cmd.Run(); err != nil {
		// If the error is ErrNotFound, we ignore it. Otherwise, return it.
		if err != exec.ErrNotFound {
			return ""
		}
	}
	passwd := strings.TrimSpace(buff.String())
	if passwd == "" {
		return ""
	}

	// username:password:uid:gid:gecos:home:shell
	passwdParts := strings.SplitN(passwd, ":", 7)
	if len(passwdParts) <= 5 {
		return ""
	}

	return passwdParts[5]
}

// windows
func dirUserProfileEnvVar() string {
	return os.Getenv("USERPROFILE")
}

// windows
func drivePathEnvVar() string {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		return ""
	}

	return home
}
