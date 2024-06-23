package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "null" {
		return fmt.Errorf("cannot unmarshal null duration for key: %s", str)
	}

	if strings.Contains(str, "d") {
		sub := str[:len(str)-1]
		value, err := strconv.Atoi(sub)
		if err != nil {
			return fmt.Errorf("cannot unmarshal duration for key: %s: %v", str, err)
		}

		d.Duration = time.Duration(value) * time.Hour * 24
		return nil
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		return fmt.Errorf("cannot unmarshal duration for key: %s: %v", str, err)
	}

	d.Duration = duration
	return nil
}
