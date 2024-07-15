package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Parse(path string, config interface{}) error {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(file).Decode(config); err != nil {
		if cErr := file.Close(); cErr != nil {
			return fmt.Errorf("config failed to parse %s: %w; config failed to close file: %v", path, err, cErr)
		}

		return fmt.Errorf("config failed to parse %s: %w", path, err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("config failed to close %s: %w", path, err)
	}

	return nil
}
