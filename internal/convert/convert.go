package convert

import (
	"encoding/json"
	"fmt"
	"log/slog"

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
