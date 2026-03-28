package dotenvsecret

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

// Load parses a .envsecret file and then loads all the variables found as environment variables.
// It uses the provided SecretManager to fetch the actual secret value.
func Load(ctx context.Context, manager SecretManager, filenames ...string) error {
	if loadDotenvsecretDisabled() {
		log.Println("dotenvsecret: .envsecret loading disabled by DOTENVSECRET_DISABLED environment variable")
		return nil
	}

	if len(filenames) == 0 {
		filenames = []string{".envsecret"}
	}

	for _, filename := range filenames {
		err := loadFile(ctx, manager, filename)
		if err != nil {
			return err
		}
	}
	return nil
}

// Unload parses a .envsecret file and then removes all the variables found from environment variables.
func Unload(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{".envsecret"}
	}

	for _, filename := range filenames {
		err := unloadFile(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadDotenvsecretDisabled() bool {
	val, ok := os.LookupEnv("DOTENVSECRET_DISABLED")
	if !ok {
		return false
	}
	val = strings.ToLower(val)
	return val == "1" || val == "true" || val == "t" || val == "yes" || val == "y"
}

func loadFile(ctx context.Context, manager SecretManager, filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envVar := strings.TrimSpace(parts[0])
			secretID := strings.TrimSpace(parts[1])

			// Optional: strip quotes if present
			if (strings.HasPrefix(secretID, "\"") && strings.HasSuffix(secretID, "\"")) ||
				(strings.HasPrefix(secretID, "'") && strings.HasSuffix(secretID, "'")) {
				secretID = secretID[1 : len(secretID)-1]
			}

			secretParts := strings.Split(secretID, ":")
			secretName := secretParts[0]
			versionID := "latest"
			if len(secretParts) == 2 {
				versionID = secretParts[1]
			}

			// In python version: version_id is "latest", source_id is None by default
			secretValue, err := manager.AccessSecret(ctx, secretName, versionID)
			if err != nil {
				fmt.Printf("Warning: Failed to load secret '%s' for environment variable '%s': %v\n", secretID, envVar, err)
				continue
			}
			os.Setenv(envVar, secretValue)
		}
	}

	return scanner.Err()
}

func unloadFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("dotenvsecret: %s file not found at %s\n", filename, filename)
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envVar := strings.TrimSpace(parts[0])
			os.Unsetenv(envVar)
		}
	}

	return scanner.Err()
}
