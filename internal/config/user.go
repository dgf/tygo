package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const (
	cfgDirName  = "tygo"
	cfgFileName = "config.json"
)

func LoadUserConfig() (Config, error) {
	var c Config

	d, err := os.UserConfigDir()
	if err != nil {
		return c, fmt.Errorf("user config dir access failed: %w", err)
	}

	b, err := os.ReadFile(path.Join(d, cfgDirName, cfgFileName))
	if err != nil {
		return c, fmt.Errorf("user config read failed: %w", err)
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		return c, fmt.Errorf("config unmarshal failed: %w", err)
	}

	return c, nil
}

func WriteUserConfig(c Config) error {
	d, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("user config dir access failed: %w", err)
	}

	dirName := path.Join(d, cfgDirName)
	dirInfo, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		err = os.Mkdir(dirName, 0o700)
		if err != nil {
			return fmt.Errorf("make user config app dir failed: %w", err)
		}

		dirInfo, err = os.Stat(dirName)
	}

	if err != nil {
		return fmt.Errorf("user config app dir access failed: %w", err)
	}

	if !dirInfo.IsDir() {
		return fmt.Errorf("user config app dir %q isn't accessible", dirName)
	}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("user config marshal failed: %w", err)
	}

	err = os.WriteFile(path.Join(dirName, cfgFileName), b, 0o600)
	if err != nil {
		return fmt.Errorf("user config write failed: %w", err)
	}

	return nil
}
