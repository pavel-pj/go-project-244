package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func parceFile(path string) (map[string]interface{}, error) {

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

func GendDiff(file01, file02, format string) (string, error) {

	//Парсим файлы
	data01, err := parceFile(file01)
	if err != nil {
		return "", err
	}
	data02, err := parceFile(file02)
	if err != nil {
		return "", err
	}

	result := []DiffItem{}

	//merge
	result1 := mergeRecursive(result, data01)
	result2 := mergeRecursive(result1, data02)
	result3 := getSorted(result2)
	//compare
	result4 := differ(result3, data01, data02)
	//format
	return formater(result4, format), nil
}

type DiffItem struct {
	key      string
	value    []interface{}
	result   string
	children []DiffItem
}

func mergeRecursive(result []DiffItem, file map[string]interface{}) []DiffItem {
	for key, value := range file {

		item := getDiffItem(result, key)
		//==============================================================================
		// Обработка простых значений
		if !isMap(value) {
			//если ключ существует, те был добавлен из первого файла, добавляем в значение  существующий слайс
			if item != nil && len(item.children) == 0 {

				//могут быть одинаковые ключи для простого и вложенного знчаения
				//existingSlice := result.value.([]interface{})
				if item.value[0] != value {
					item.value = append(item.value, saveCorrectValues(value))
				}
				// Если значения одинаковые - оставляем срез как есть
				continue
			}

			// Ключа нет в результате - создаем срез с одним значением
			result = append(result, DiffItem{
				key:      key,
				value:    []interface{}{saveCorrectValues(value)},
				result:   "",
				children: []DiffItem{},
			})

			continue

		}
		//==============================================================================
		// проверяем вложенные данные

		nestedMap := value.(map[string]interface{})
		if item != nil && len(item.children) > 0 {
			// Если такой ключ с вложенным значением уже существуе
			item.children = mergeRecursive(item.children, nestedMap)
			continue
		}

		//Если папка не существует, создаем ее
		nestedChilds := mergeRecursive([]DiffItem{}, nestedMap)
		result = append(result, DiffItem{
			key:      key,
			value:    []interface{}{},
			result:   "",
			children: nestedChilds,
		})

	}

	return result
}

func getDiffItem(result []DiffItem, key string) *DiffItem {
	for i := range result {
		if result[i].key == key {
			return &result[i]
		}
	}
	return nil
}

func isMap(value interface{}) bool {

	if _, ok := value.(map[string]interface{}); ok {
		return true
	}
	return false
}

func saveCorrectValues(value interface{}) interface{} {
	if value == nil {
		return "null"
	}
	return value
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

func formater(diff []DiffItem, format string) string {

	level := 0
	switch format {
	case "AAAA":

		return "{\n" + recursiveFormat(diff, level) + "}\n"
	default:
		return "{\n" + recursiveFormat(diff, level) + "}\n"

	}

}

func recursiveFormat(diff []DiffItem, curLevel int) string {

	step := 4
	smb := " "
	indent := ""
	result := ""
	for _, r := range diff {
		//для простых значений
		if len(r.value) == 1 {
			indent = strings.Repeat(smb, curLevel*step) + getSymbol(r)
			result += indent + r.key + ": " + getValue(r.value[0]) + "\n"
			continue
		}
		if len(r.value) == 2 {
			indent = strings.Repeat(smb, curLevel*step) + getSymbol(r)
			result += indent + r.key + ": " + getValue(r.value[0]) + "\n"
			continue
		}
		if len(r.children) > 0 {
			//Для вложенных
			indent = strings.Repeat(smb, curLevel*step)

			result += indent + getSymbol(r) + r.key + ": {\n"
			result += recursiveFormat(r.children, curLevel+1)
			result += indent + strings.Repeat(smb, 2) + "}\n"
		}

	}
	return result
}

func getSymbol(item DiffItem) string {
	switch item.result {
	case "new":
		return "+ "
	case "deleted":
		return "- "
	case "unchanged":
		return "  "
	}
	return "  "
}

func getValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}
