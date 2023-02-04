package security

import (
	"os/user"

	"github.com/luisnquin/nao/v3/internal"
	"github.com/zalando/go-keyring"
)

// Retrieves a secret stored in the system keyring. It uses the current
// username of the caller and the app name to access the secret.
//
// If the current username can't be determined, the function panics
// with an error.
func GetSecretFromKeyring() (string, error) {
	caller, err := user.Current()
	if err != nil {
		panic(err)
	}

	return keyring.Get(internal.AppName, caller.Username)
}

// Stores a secret in the system keyring. It uses the current username
// of the caller and the app name to access the secret.
//
// If the current username can't be determined, the function panics
// with an error.
func SetSecretInKeyring(secret string) error {
	caller, err := user.Current()
	if err != nil {
		panic(err)
	}

	return keyring.Set(internal.AppName, caller.Username, secret)
}
