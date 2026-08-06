package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Парсит файл json/yaml в map новый
func ParseFile(path string) (map[string]interface{}, error) {

	// #nosec G304 - path comes from trusted source (config file)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	ext := filepath.Ext(path)
	fileType := strings.TrimPrefix(ext, ".")

	var result map[string]interface{}

	switch fileType {
	case "json":
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("not a JSON array or parsing error: %w", err)
		}
	case "yml":
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("not a YAML array or parsing error: %w", err)
		}
	case "yaml":
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("not a YAML array or parsing error: %w", err)
		}

	}

	return result, nil

}
