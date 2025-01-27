package user

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// UsernameCache controls caching of the username. Caching is enabled
// by default.
var UsernameCache = true

var (
	usernameCacheValue string
	usernameCacheLock  sync.RWMutex
)

// Reset clears the cache, forcing the next call to Dir to re-detect
// the home directory. This generally never has to be called, but can be
// useful in tests if you're modifying the home directory via the HOME
// env var or something.
func ResetUsername() {
	usernameCacheLock.Lock()
	defer usernameCacheLock.Unlock()
	usernameCacheValue = ""
}

// Username returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Username() (value string, fault error) {
	if UsernameCache {
		usernameCacheLock.RLock()
		cached := usernameCacheValue
		usernameCacheLock.RUnlock()
		if cached != "" {
			return cached, nil
		}
	}

	usernameCacheLock.Lock()
	defer usernameCacheLock.Unlock()

	result, err := uname()
	if err != nil {
		return "", err
	}
	usernameCacheValue = result
	return result, nil
}

var usrMethods = []func() string{}

func uname() (value string, fault error) {
	for _, m := range usrMethods {
		if usr := m(); usr != "" {
			return usr, nil
		}
	}

	return "", errors.New("could not determine username")
}

var userEnv = "USER"

// OS agnostic, base-level method to get the home directory. This is the
// first method to be tried. The value for homeEnv is determined in the
// OS-specific files.
func userEnvVar() string {
	return os.Getenv(userEnv)
}

// Mac OS/Darwin, Linux. Plan9
func usrWhoAmI() string {
	var buff bytes.Buffer
	cmd := exec.Command("whoami")
	cmd.Stdout = &buff
	if err := cmd.Run(); err == nil {
		result := strings.TrimSpace(buff.String())
		if result != "" {
			return result
		}
	}
	return ""
}

// plan9, linux
func usrGetEntPasswd() string {
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

	return passwdParts[0]
}
