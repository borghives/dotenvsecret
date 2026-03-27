package dotenvsecret

import (
	"context"
)

// SecretManager is an interface that allows fetching secrets from different backends.
type SecretManager interface {
	AccessSecret(ctx context.Context, secretID, versionID string) (string, error)
}
