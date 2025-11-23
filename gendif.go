package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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
	case "yaml":
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
	result1 := mergeRecursive(result, data01, "")
	result2 := mergeRecursive(result1, data02, "")

	//fmt.Println(result2)
	//return "", nil

	result3 := getSorted(result2)
	//compare
	result4 := differ(result3, data01, data02)
	fmt.Println(result4)
	fmt.Println(" ")
	//format
	return formater(result4, format), nil

}

type DiffItem struct {
	key      string
	value    []interface{}
	result   string
	children []DiffItem
	path     string
}

func mergeRecursive(result []DiffItem, file map[string]interface{}, path string) []DiffItem {
	for key, value := range file {
		curPath := key
		if len(path) > 0 {
			curPath = path + "." + key
		}
		item := getDiffItem(result, key)
		//==============================================================================
		// Обработка простых значений
		if !isMap(value) {

			//1.
			// если ключ существует, НО был стуктурой
			if item != nil && len(item.children) > 0 {

				item.value = append(item.value, item.children)
				item.value = append(item.value, value)
				//удалим стукруту с chilld , она уже не потребуется, все значения old & new в слайсе
				item.children = []DiffItem{}
				continue
			}
			//2.
			// если ключ существует, НО c плоской стуктурой
			if item != nil && len(item.children) == 0 && item.value[0] != value {
				item.value = append(item.value, value)
				continue
			}
			//3.
			// если ключ существует и он равен текущему
			if item != nil && item.value[0] == value {
				continue
			}

			//4.
			//Ключа нет в результате - создаем срез с одним значением
			result = append(result, DiffItem{
				key:      key,
				value:    []interface{}{value},
				result:   "",
				children: []DiffItem{},
				path:     curPath,
			})

			continue
		}
		//==============================================================================
		// проверяем вложенные данные

		nestedMap := value.(map[string]interface{})
		//1.
		// Если такой ключ с вложенным значением уже существует
		if item != nil && len(item.children) > 0 {
			item.children = mergeRecursive(item.children, nestedMap, curPath)
			continue
		}

		//получаем вложенные папки
		nestedChilds := mergeRecursive([]DiffItem{}, nestedMap, curPath)

		//2.
		//Если такой ключ существует, но значение - простое
		if item != nil && len(item.children) == 0 && len(item.value) > 0 {
			item.value = append(item.value, nestedChilds)
			continue
		}

		//Если папка не существует, создаем ее
		result = append(result, DiffItem{
			key:      key,
			value:    []interface{}{},
			result:   "",
			children: nestedChilds,
			path:     curPath,
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

/*
func saveCorrectValues(value interface{}) interface{} {
	if value == nil {
		return "null"
	}
	return value
}
*/

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
	case "plain":

		return formatPlain(diff)
	default:
		return "{\n" + formatStylish(diff, level) + "}\n"

	}

}

func formatStylish(diff []DiffItem, curLevel int) string {

	step := 4
	smb := " "
	indent := ""
	result := ""
	for _, r := range diff {
		//1
		//Простые значений Добавлено/Удалено
		if len(r.value) == 1 {
			indent = strings.Repeat(smb, curLevel*step) + getSymbol(r)
			result += indent + r.key + ": " + getValue(r.value[0]) + "\n"
			continue
		}

		//2.
		//Обновлено для простых и рекурсивных значений
		if len(r.value) == 2 {
			//result += "ВЛОЖЕННАЯ \n"

			if reflect.TypeOf(r.value[0]) != reflect.TypeOf([]DiffItem{}) {
				indent = strings.Repeat(smb, curLevel*step) + "- "

				result += indent + r.key + ": " + getValue(r.value[0]) + "\n"

			}
			if reflect.TypeOf(r.value[1]) != reflect.TypeOf([]DiffItem{}) {

				indent = strings.Repeat(smb, curLevel*step) + "+ "

				result += indent + r.key + ": " + getValue(r.value[1]) + "\n"
			}

			//indent = strings.Repeat(smb, curLevel*step) + getSymbol(r)
			//result += indent + r.key + ": " + getValue(r.value[0]) + "\n"
			continue
		}
		//3 для вложенных элементов
		if len(r.children) > 0 {
			//Для вложенных
			indent = strings.Repeat(smb, curLevel*step)

			result += indent + getSymbol(r) + r.key + ": {\n"
			result += formatStylish(r.children, curLevel+1)
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

func formatPlain(diff []DiffItem) string {
	result := ""
	/*
		for _, r := range diff {

			switch r.result {
			case "deleted":
				result += "Property '" + r.path + "' was removed\n"
			case "new":
				value := ""
				if len(r.children) > 0 {
					value1 = "[complex value]"
				} else {
					value1 = getValue(r.value[0])
				}
				result += "Property '" + r.path + "' was added with value: " + value + "\n"
			case "updated":
				if len(r.children) == 0 {
					result += "Property '" + r.path + "' was updated. From " + r.value[0] + " to " + r.value[1] + "\n"
					continue
				}
				result += "Property '" + r.path + "' was updated. From " + r.value[0] + " to " + r.value[1] + "\n"

			default:
				result += ""
			}

			if len(r.children) > 0 {
				result += formatPlain(r.children)
			}

		}*/
	return result
}
