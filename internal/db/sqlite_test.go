package db

import (
	"testing"
)

var (
	un string = "username"
	pw string = "password"
)

func TestSignupAndLogin(t *testing.T) {
	var err error

	err = InitTestDb()
	if err != nil {
		t.Error(err)
		return
	}
	defer CleanupTestDb()

	err = SignupUser(un, pw, "user")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = LoginUser(un, pw)
	if err != nil {
		t.Error(err)
		return
	}
}
