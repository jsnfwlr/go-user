//go:build windows

package user

func init() {
	dirMethods = []func() string{
		homeEnvVar,
		dirUserProfileEnvVar,
		drivePathEnvVar,
	}

	userEnv = "USERNAME"

	usrMethods = []func() string{
		usrWhoAmI,
	}
}
