package configuration

import "flag"

func Load[TConfig any]() (TConfig, error) {
	path := flag.String("c", "config.json", "configuration file path")
	flag.Parse()

	var cfg TConfig
	if err := parseJSON(*path, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
