package user

import (
	"os/user"
	"testing"
)

func TestUsername(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usr, err := Username()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if u.Username != usr {
		t.Fatalf("%#v != %#v", u.Username, usr)
	}

	UsernameCache = true
	defer func() { UsernameCache = false }()
	defer patchEnv(userEnv, "")()
	usr, err = Username()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if u.Username != usr {
		t.Fatalf("%#v != %#v", u.Username, usr)
	}
}
