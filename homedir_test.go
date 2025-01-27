package user

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func patchEnv(key, value string) func() {
	bck := os.Getenv(key)
	deferFunc := func() {
		os.Setenv(key, bck)
	}

	if value != "" {
		os.Setenv(key, value)
	} else {
		os.Unsetenv(key)
	}

	return deferFunc
}

func BenchmarkDir(b *testing.B) {
	// We do this for any "warmups"
	for i := 0; i < 10; i++ {
		_, _ = HomeDir()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HomeDir()
	}
}

func TestHomeDir(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	dir, err := HomeDir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if u.HomeDir != dir {
		t.Fatalf("%#v != %#v", u.HomeDir, dir)
	}

	HomeDirCache = true
	defer func() { HomeDirCache = false }()
	defer patchEnv(homeEnv, "")()
	dir, err = HomeDir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if u.HomeDir != dir {
		t.Fatalf("%#v != %#v", u.HomeDir, dir)
	}
}

func TestExpandPath(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	cases := []struct {
		Input  string
		Output string
		Err    bool
	}{
		{
			"/foo",
			"/foo",
			false,
		},

		{
			"~/foo",
			filepath.Join(u.HomeDir, "foo"),
			false,
		},

		{
			"",
			"",
			false,
		},

		{
			"~",
			u.HomeDir,
			false,
		},

		{
			"~/foo/../foo",
			filepath.Join(u.HomeDir, "foo"),
			false,
		},

		{
			"~foo/foo",
			"",
			true,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s-%s", tc.Input, tc.Output), func(t *testing.T) {
			actual, err := ExpandPath(tc.Input)
			if (err != nil) != tc.Err {
				t.Fatalf("Input: %#v\n\nErr: %s", tc.Input, err)
			}

			if actual != tc.Output {
				t.Fatalf("Input: %#v\n\nOutput: %#v", tc.Input, actual)
			}
		})
	}

	t.Run("custom no cache", func(t *testing.T) {
		HomeDirCache = false
		defer func() { HomeDirCache = true }()
		defer patchEnv("HOME", "/custom/path/")()
		expected := filepath.Join("/", "custom", "path", "foo/bar")
		actual, err := ExpandPath("~/foo/bar")

		if err != nil {
			t.Errorf("No error is expected, got: %v", err)
		} else if actual != expected {
			t.Errorf("Expected: %v; actual: %v", expected, actual)
		}
	})
}
