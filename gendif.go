package code

import (
	"code/formatters"
	"code/parser"
	"code/types"
	"sort"
	"strings"
)

func GenDiff(file01, file02, format string) (string, error) {

	//Парсим файлы
	data01, err := parser.ParceFile(file01)
	if err != nil {
		return "", err
	}
	data02, err := parser.ParceFile(file02)
	if err != nil {
		return "", err
	}

	result := []types.DiffItem{}

	//merge
	result1 := mergeRecursive(result, data01, "")
	result2 := mergeRecursive(result1, data02, "")

	result3 := getSorted(result2)
	//compare
	result4 := differ(result3, data01, data02)
	//format
	return formatters.Formater(result4, format), nil

}

func mergeRecursive(result []types.DiffItem, file map[string]interface{}, path string) []types.DiffItem {
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
			if item != nil && len(item.Children) > 0 {

				item.Value = append(item.Value, item.Children)
				item.Value = append(item.Value, value)
				//удалим стукруту с chilld , она уже не потребуется, все значения old & new в слайсе
				item.Children = []types.DiffItem{}
				continue
			}
			//2.
			// если ключ существует, НО c плоской стуктурой
			if item != nil && len(item.Children) == 0 && item.Value[0] != value {
				item.Value = append(item.Value, value)
				continue
			}
			//3.
			// если ключ существует и он равен текущему
			if item != nil && item.Value[0] == value {
				continue
			}

			//4.
			//Ключа нет в результате - создаем срез с одним значением
			result = append(result, types.DiffItem{
				Key:      key,
				Value:    []interface{}{value},
				Result:   "",
				Children: []types.DiffItem{},
				Path:     curPath,
			})

			continue
		}
		//==============================================================================
		// проверяем вложенные данные

		nestedMap := value.(map[string]interface{})
		//1.
		// Если такой ключ с вложенным значением уже существует
		if item != nil && len(item.Children) > 0 {
			item.Children = mergeRecursive(item.Children, nestedMap, curPath)
			continue
		}

		//получаем вложенные папки
		nestedChilds := mergeRecursive([]types.DiffItem{}, nestedMap, curPath)

		//2.
		//Если такой ключ существует, но значение - простое
		if item != nil && len(item.Children) == 0 && len(item.Value) > 0 {
			item.Value = append(item.Value, nestedChilds)
			continue
		}

		//Если папка не существует, создаем ее
		result = append(result, types.DiffItem{
			Key:      key,
			Value:    []interface{}{},
			Result:   "",
			Children: nestedChilds,
			Path:     curPath,
		})

	}

	return result
}

func getDiffItem(result []types.DiffItem, key string) *types.DiffItem {
	for i := range result {
		if result[i].Key == key {

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

func getSorted(diff []types.DiffItem) []types.DiffItem {

	// Сортируем текущий уровень
	sort.Slice(diff, func(i, j int) bool {
		return strings.ToLower(diff[i].Key) < strings.ToLower(diff[j].Key)
	})

	// Рекурсивно сортируем детей для каждого элемента
	for i := range diff {
		if len(diff[i].Children) > 0 {
			diff[i].Children = getSorted(diff[i].Children)
		}
	}

	return diff
}

func differ(diff []types.DiffItem, file01 map[string]interface{}, file02 map[string]interface{}) []types.DiffItem {

	for i := range diff {

		fileChild01, inFile01 := file01[diff[i].Key]
		fileChild02, inFile02 := file02[diff[i].Key]

		// Проверяем что оба значения - map[string]interface{}
		childMap01, ok01 := fileChild01.(map[string]interface{})
		childMap02, ok02 := fileChild02.(map[string]interface{})

		//Для вложенных структур
		if len(diff[i].Children) > 0 {

			if !inFile01 && inFile02 {
				diff[i].Result = "new"
				continue
			} else if inFile01 && !inFile02 {
				diff[i].Result = "deleted"
				continue
			}

			if ok01 && ok02 {
				diff[i].Children = differ(diff[i].Children, childMap01, childMap02)
			}
			continue
		}

		//Для конечных нод
		if len(diff[i].Value) == 2 {
			diff[i].Result = "updated"
			continue
		}
		if !inFile01 && inFile02 {
			diff[i].Result = "new"
			continue
		} else if inFile01 && !inFile02 {
			diff[i].Result = "deleted"
			continue
		}

		diff[i].Result = "unchanged"

	}

	return diff

}
