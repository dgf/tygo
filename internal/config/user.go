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

type UserDirAccessError struct {
	Dir string
}

func (e *UserDirAccessError) Error() string {
	return fmt.Sprintf("user config app dir %q isn't accessible", e.Dir)
}

func LoadUserConfig() (Config, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Default(), fmt.Errorf("user config dir access failed: %w", err)
	}

	root, err := os.OpenRoot(dir)
	if err != nil {
		return Default(), fmt.Errorf("user config dir open failed: %w", err)
	}
	defer func() {
		_ = root.Close()
	}()

	data, err := root.ReadFile(path.Join(cfgDirName, cfgFileName))
	if err != nil {
		return Default(), fmt.Errorf("user config read failed: %w", err)
	}

	var v Config

	err = json.Unmarshal(data, &v)
	if err != nil {
		return v, fmt.Errorf("config unmarshal failed: %w", err)
	}

	return v, nil
}

func WriteUserConfig(cfg Config) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("user config dir access failed: %w", err)
	}

	dirName := path.Join(dir, cfgDirName)
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
		return &UserDirAccessError{Dir: dirName}
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("user config marshal failed: %w", err)
	}

	err = os.WriteFile(path.Join(dirName, cfgFileName), b, 0o600)
	if err != nil {
		return fmt.Errorf("user config write failed: %w", err)
	}

	return nil
}
