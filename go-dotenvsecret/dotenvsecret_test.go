package dotenvsecret

import (
	"context"
	"os"
	"testing"
)

type MockSecretManager struct {
	secrets map[string]string
}

func (m *MockSecretManager) AccessSecret(ctx context.Context, secretID, versionID string) (string, error) {
	if val, ok := m.secrets[secretID]; ok {
		return val, nil
	}
	return "", os.ErrNotExist
}

func TestLoadDotenvsecret(t *testing.T) {
	// Create a temporary .envsecret file
	content := `
# A comment
DB_PASS = "my_database_password"
API_KEY=my_api_key
`
	err := os.WriteFile(".envsecret.test", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	defer os.Remove(".envsecret.test")

	manager := &MockSecretManager{
		secrets: map[string]string{
			"my_database_password": "supersecretpassword",
			"my_api_key":           "1234567890abcdef",
		},
	}

	err = Load(context.Background(), manager, ".envsecret.test")
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if val := os.Getenv("DB_PASS"); val != "supersecretpassword" {
		t.Errorf("Expected DB_PASS to be 'supersecretpassword', got '%s'", val)
	}

	if val := os.Getenv("API_KEY"); val != "1234567890abcdef" {
		t.Errorf("Expected API_KEY to be '1234567890abcdef', got '%s'", val)
	}

	err = Unload(".envsecret.test")
	if err != nil {
		t.Fatalf("Failed to unload: %v", err)
	}

	if val := os.Getenv("DB_PASS"); val != "" {
		t.Errorf("Expected DB_PASS to be empty, got '%s'", val)
	}
	if val := os.Getenv("API_KEY"); val != "" {
		t.Errorf("Expected API_KEY to be empty, got '%s'", val)
	}
}

func TestLoadDisabled(t *testing.T) {
	os.Setenv("DOTENVSECRET_DISABLED", "1")
	defer os.Unsetenv("DOTENVSECRET_DISABLED")

	err := os.WriteFile(".envsecret.test", []byte("VAR=secret"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	defer os.Remove(".envsecret.test")

	manager := &MockSecretManager{
		secrets: map[string]string{"secret": "val"},
	}

	err = Load(context.Background(), manager, ".envsecret.test")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	if os.Getenv("VAR") == "val" {
		t.Errorf("Expected VAR environment variable to not be set")
	}
}
