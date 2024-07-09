package cache

import (
	"encoding/json"
	"os"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

func readFromFile[TMessage any](path string, logger log.Logger) ([]TMessage, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cErr := file.Close(); cErr != nil {
			logger.Error().Msgf("failed close logs file: %v", err)
		}
	}()

	var result []TMessage
	if err = json.NewDecoder(file).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func saveToFile[TMessage any](path string, result []TMessage, logger log.Logger) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := file.Close(); cErr != nil {
			logger.Error().Msgf("failed close logs file: %v", err)
		}
	}()

	if err = json.NewEncoder(file).Encode(result); err != nil {
		return err
	}

	return nil
}

func removeFile(path string, logger log.Logger) {
	if err := os.Remove(path); err != nil {
		logger.Error().Msgf("failed to remove logs file: %v", err)
	}
}
