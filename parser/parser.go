package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParceFile(path string) (map[string]interface{}, error) {

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//сначала проверяется относительынй путь в текущей папке
	var data []byte
	data, err = os.ReadFile(wd + "/" + path)
	if err != nil {

		//Абсолютный путь
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", path, err)
		}
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
