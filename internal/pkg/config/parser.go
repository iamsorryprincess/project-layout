package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const (
	envDev   = "dev"
	envStage = "stage"
	envProd  = "prod"
)

func formatPath(env string, serviceName string) string {
	return fmt.Sprintf("configs/%s/%s.config.json", env, serviceName)
}

func Parse(serviceName string, config interface{}) error {
	env := flag.String("e", envDev, "service environment")
	flag.Parse()

	path := ""
	switch *env {
	case envDev:
		path = formatPath(envDev, serviceName)
	case envStage:
		path = formatPath(envStage, serviceName)
	case envProd:
		path = formatPath(envProd, serviceName)
	default:
		return fmt.Errorf("unknown environment: %s", *env)
	}

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
