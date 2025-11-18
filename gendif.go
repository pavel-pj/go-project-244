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
	key      string
	value    []interface{}
	result   string
	children []DiffItem
}

func GendDiff02(file01, file02 map[string]interface{}) []DiffItem {

	result := make(map[string]interface{})

	//merge
	result1 := mergeRecursive(result, file01)
	result2 := mergeRecursive(result1, file02)
	result3 := getUsefulFormat(result2)
	result4 := getSorted(result3)

	result5 := differ(result4, file01, file02)

	return result5
}

/*
	func mergeRecursive(result map[string]interface{}, file map[string]interface{}) map[string]interface{} {
		for key, value := range file {
			if !isMap(value) {

				// Обработка простых значений
				if existedValue, ok := result[key]; ok {
					// Если ключ уже существует - добавляем значение в срез
					existingSlice := existedValue.([]interface{})
					if existingSlice[0] != value {
						result[key] = append(existingSlice, value)
					}
					// Если значения одинаковые - оставляем срез как есть
				} else {
					// Ключа нет в результате - создаем срез с одним значением
					result[key] = []interface{}{value}
				}

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
*/
func mergeRecursive(result map[string]interface{}, file map[string]interface{}) map[string]interface{} {
	for key, value := range file {
		if !isMap(value) {
			// Обработка простых значений
			if existedValue, ok := result[key]; ok {
				// Проверяем тип существующего значения
				switch existing := existedValue.(type) {
				case []interface{}:
					// Если уже есть срез - добавляем значение если отличается
					if existing[0] != value {
						result[key] = append(existing, value)
					}
				case map[string]interface{}:
					// Если уже есть map - преобразуем в срез
					result[key] = []interface{}{existing, value}
				default:
					// Если одиночное значение - создаем срез
					if existing != value {
						result[key] = []interface{}{existing, value}
					}
				}
			} else {
				// Ключа нет в результате - создаем срез с одним значением
				result[key] = []interface{}{value}
			}
		} else {
			// Handle nested map
			if existing, exists := result[key]; exists {
				// Проверяем тип существующего значения
				switch existingMap := existing.(type) {
				case map[string]interface{}:
					// Recursively merge into existing map
					mergeRecursive(existingMap, value.(map[string]interface{}))
				case []interface{}:
					// Если уже есть срез - нужно проверить можно ли мерджить
					// Для простоты создаем новую map и мерджим оба значения
					nestedResult := make(map[string]interface{})
					// Пытаемся извлечь map из среза и мерджить
					if len(existingMap) > 0 {
						if firstMap, ok := existingMap[0].(map[string]interface{}); ok {
							mergeRecursive(nestedResult, firstMap)
						}
					}
					mergeRecursive(nestedResult, value.(map[string]interface{}))
					result[key] = []interface{}{nestedResult}
				default:
					// Если одиночное значение - создаем новую map
					nestedResult := make(map[string]interface{})
					mergeRecursive(nestedResult, value.(map[string]interface{}))
					result[key] = nestedResult
				}
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

func getUsefulFormat(rawArray map[string]interface{}) []DiffItem {
	sortedDiff := []DiffItem{}

	for key, value := range rawArray {
		if slice, ok := value.([]interface{}); ok {
			// Простой случай - массив значений
			sortedDiff = append(sortedDiff, DiffItem{
				key:    key,
				value:  slice,
				result: "",
			})
		} else {
			// Рекурсивный случай - вложенная map
			nestedMap := value.(map[string]interface{})
			nestedItems := getUsefulFormat(nestedMap)
			sortedDiff = append(sortedDiff, DiffItem{
				key:      key,
				value:    nil, // или можно сделать []interface{}{}
				result:   "",
				children: nestedItems, // сохраняем вложенные элементы
			})
		}
	}

	return sortedDiff
}

func getSorted(diff []DiffItem) []DiffItem {

	// Сортируем текущий уровень
	sort.Slice(diff, func(i, j int) bool {
		return strings.ToLower(diff[i].key) < strings.ToLower(diff[j].key)
	})

	// Рекурсивно сортируем детей для каждого элемента
	for i := range diff {
		if len(diff[i].children) > 0 {
			diff[i].children = getSorted(diff[i].children)
		}
	}

	return diff
}

func differ(diff []DiffItem, file01 map[string]interface{}, file02 map[string]interface{}) []DiffItem {

	for i := range diff {

		fileChild01, inFile01 := file01[diff[i].key]
		fileChild02, inFile02 := file02[diff[i].key]

		// Проверяем что оба значения - map[string]interface{}
		childMap01, ok01 := fileChild01.(map[string]interface{})
		childMap02, ok02 := fileChild02.(map[string]interface{})

		//Для вложенных структур
		if len(diff[i].children) > 0 {

			if !inFile01 && inFile02 {
				diff[i].result = "new"
				continue
			} else if inFile01 && !inFile02 {
				diff[i].result = "deleted"
				continue
			}

			if ok01 && ok02 {
				diff[i].children = differ(diff[i].children, childMap01, childMap02)
			}
			continue
		}

		//Для конечных нод
		if len(diff[i].value) == 2 {
			diff[i].result = "updated"
			continue
		}
		if !inFile01 && inFile02 {
			diff[i].result = "new"
			continue
		} else if inFile01 && !inFile02 {
			diff[i].result = "deleted"
			continue
		}

		diff[i].result = "unchanged"

	}

	return diff

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
