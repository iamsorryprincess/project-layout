package cfg

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

func trim(line string) string {
	return strings.Trim(strings.TrimSpace(line), `{,"}`)
}

func Parse(path string, config interface{}) error {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(file).Decode(config); err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	keys := make(map[string]string)

	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return err
		}

		line := trim(scanner.Text())

		if !strings.Contains(line, ":") {
			continue
		}

		before, after, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		before = trim(before)
		after = trim(after)

		keys[before] = after
	}

	return nil
}
