package config_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dgf/tygo/internal/config"
)

const lastWorkingConfigExample = `{
  "version": 1,
  "dict": "german",
  "strict": false,
  "top": 100,
  "count": 20,
  "width": 30,
  "nums": true,
  "punct": true,
  "freqs": {
    "word": 85,
    "number": 7,
    "period": 12,
    "comma": 8,
    "quotation": 3,
    "question": 4,
    "exclamation": 3,
    "brackets": 2,
    "braces": 2,
    "parenthesis": 3,
    "colon": 3,
    "semicolon": 2
  }
}`

const nextSavedConfigExample = `{
  "version": 2,
  "dict": "german",
  "strict": false,
  "top": 100,
  "count": 20,
  "width": 30,
  "nums": true,
  "punct": true,
  "noRepeat": 5,
  "freqs": {
    "word": 85,
    "number": 7,
    "period": 12,
    "comma": 8,
    "quotation": 3,
    "question": 4,
    "exclamation": 3,
    "brackets": 2,
    "braces": 2,
    "parenthesis": 3,
    "colon": 3,
    "semicolon": 2
  }
}`

// This test must be adjusted whenever the config changes,
// that's on purpose, to simulate the last migration step(s).
func TestMigrate(t *testing.T) {
	t.Parallel()

	var cfg config.Config

	err := json.Unmarshal([]byte(lastWorkingConfigExample), &cfg)
	if err != nil {
		t.Error(err)
	}

	migrated := config.Migrate(&cfg)
	if migrated == false {
		t.Error("invalid migration test setup")
	}

	migratedConfig, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Error(err)
	}

	if strings.Compare(string(migratedConfig), nextSavedConfigExample) != 0 {
		t.Errorf("invalid migration result, got:\n%v\n\nwant:\n%v", string(migratedConfig), nextSavedConfigExample)
	}
}

func TestDefaultVersion(t *testing.T) {
	t.Parallel()

	defaultVersion := config.Default().Version
	migrateVersion := len(config.Migrations())

	if defaultVersion != migrateVersion {
		t.Errorf("Update default version, want: %d, got: %d", migrateVersion, defaultVersion)
	}
}
