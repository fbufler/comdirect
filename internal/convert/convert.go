package convert

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"gopkg.in/yaml.v2"
)

func JSONToYAML(data string) (string, error) {
	inp := []byte(data)
	slog.Debug(string(inp))

	var parsed interface{}
	if err := json.Unmarshal(inp, &parsed); err != nil {
		return "", err
	}
	slog.Debug(fmt.Sprint(parsed))

	out, err := yaml.Marshal(parsed)
	if err != nil {
		return "", err
	}
	slog.Debug(string(out))

	return string(out), nil
}

func JSONToReadableJSON(data string) (string, error) {
	inp := []byte(data)
	slog.Debug(string(inp))

	var parsed interface{}
	if err := json.Unmarshal(inp, &parsed); err != nil {
		return "", err
	}
	slog.Debug(fmt.Sprint(parsed))

	out, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return "", err
	}
	slog.Debug(string(out))

	return string(out), nil
}

func TimeStringToTime(data string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
		"02.01.2006",
		"02.01.06",
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, data)
		if err == nil {
			break
		}
	}
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
