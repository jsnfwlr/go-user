//go:build linux

package user

func init() {
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
