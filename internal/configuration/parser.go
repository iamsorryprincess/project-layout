package configuration

import (
	"errors"
	"flag"
	"fmt"
)

const TypeJSON = "json"

var ErrNotImplementedConfigType = errors.New("not implemented config type")

func Parse(configType string, serviceName string, config interface{}) error {
	switch configType {
	case TypeJSON:
		path := flag.String("c", fmt.Sprintf("configs/local/%s.config.json", serviceName), "config path")
		flag.Parse()
		return ParseJSON(*path, config)
	default:
		return ErrNotImplementedConfigType
	}
}
