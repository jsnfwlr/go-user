//go:build plan9

package user

func init() {
	// On plan9, env vars are lowercase.
	homeEnv = "home"
	userEnv = "user"

	dirMethods = []func() string{
		homeEnvVar,
		dirGetEntPasswd,
		cdPwd,
	}

	usrMethods = []func() string{
		userEnvVar,
		usrGetEntPasswd,
		usrWhoAmI,
	}
}
