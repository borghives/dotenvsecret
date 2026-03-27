package dotenvsecret

import (
	"context"
	"os"
	"os/user"

	"github.com/zalando/go-keyring"
)

type LocalKeyring struct{}

func NewLocalKeyring() *LocalKeyring {
	return &LocalKeyring{}
}

func (m *LocalKeyring) AccessSecret(ctx context.Context, secretID, versionID string) (string, error) {
	username := os.Getenv("LOCAL_KEYRING_USERNAME")
	if username == "" {
		u, err := user.Current()
		if err == nil {
			username = u.Username
		} else {
			username = "default" // fallback if user.Current() fails
		}
	}

	secret, err := keyring.Get(secretID, username)
	if err != nil {
		return "", err
	}
	return secret, nil
}
