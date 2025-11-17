package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

	}

	return result, nil

}

type DiffItem struct {
	key    string
	value  []interface{}
	result string
}

func GendDiff02(file01, file02 map[string]interface{}) map[string]interface{} {

	result := make(map[string]interface{})

	//merge
	result2 := mergeRecursive(result, file01)
	return result2
}

func mergeRecursive(result map[string]interface{}, file map[string]interface{}) map[string]interface{} {
	for key, value := range file {
		if !isMap(value) {
			result[key] = value
		} else {
			// Handle nested map
			if existing, exists := result[key]; exists && isMap(existing) {
				// Recursively merge into existing map
				existingMap := existing.(map[string]interface{})
				mergeRecursive(existingMap, value.(map[string]interface{}))
			} else {
				// Create new nested map
				nestedResult := make(map[string]interface{})
				mergeRecursive(nestedResult, value.(map[string]interface{}))
				result[key] = nestedResult
			}
		}
	}
	return result
}

func isMap(value interface{}) bool {

	if _, ok := value.(map[string]interface{}); ok {
		return true
	}
	return false
}

func GendDiff01(file01, file02 map[string]interface{}) []DiffItem {
	result := make(map[string][]interface{})

	//merge
	for k, v := range file01 {
		result[k] = append(result[k], v)
	}
	for k, v := range file02 {
		if value, ok := result[k]; ok {
			if value[0] != v {
				result[k] = append(result[k], v)
			}
		} else {
			result[k] = append(result[k], v)
		}
	}

	sortedDiff := make([]DiffItem, 0, len(result))
	for k, v := range result {
		sortedDiff = append(sortedDiff, DiffItem{key: k, value: v, result: ""})
	}
	//Сравнение

	for i := range sortedDiff {

		r := &sortedDiff[i]

		_, inFile01 := file01[r.key]
		_, inFile02 := file02[r.key]

		switch {
		case inFile01 && inFile02:
			//Изменено значение
			if len(r.value) > 1 {
				r.result = "updated"
			} else {
				r.result = "unchanged"
			}
		case !inFile01 && inFile02:
			r.result = "deleted"

		case inFile01 && !inFile02:
			r.result = "added"
		}
	}

	//Сортировка
	sort.Slice(sortedDiff, func(i, j int) bool {
		return sortedDiff[i].key < sortedDiff[j].key
	})
	return sortedDiff

	//for _, item := range sortedDiff {
	//	fmt.Printf("%s: %v -> %s\n", item.key, item.value, item.result)
	//}

}

func Format(diff []DiffItem) string {

	indent := 2
	symbol := " "
	result := ""
	for _, r := range diff {

		switch r.result {
		case "unchanged":
			result += strings.Repeat(symbol, indent*2)
			result += fmt.Sprintf("%s: %v\n", r.key, r.value)
		case "added":
			result += strings.Repeat(symbol, indent) + "+ "
			result += fmt.Sprintf("%s: %v\n", r.key, r.value)
		case "deleted":
			result += strings.Repeat(symbol, indent) + "- "
			result += fmt.Sprintf("%s: %v\n", r.key, r.value)
		case "updated":
			result += strings.Repeat(symbol, indent) + "- "
			result += fmt.Sprintf("%s: %v\n", r.key, r.value[0])
			result += strings.Repeat(symbol, indent) + "+ "
			result += fmt.Sprintf("%s: %v\n", r.key, r.value[1])

		}

	}
	return result

}
