//go:build darwin

package user

func init() {
	dirMethods = []func() string{
		homeEnvVar,
		dirDSCL,
		cdPwd,
	}

	usrMethods = []func() string{
		userEnvVar,
		usrWhoAmI,
	}
}
