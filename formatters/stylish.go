package formatters

import (
	"code/types"
	"strings"
)

func formatStylish(diff []types.DiffItem, curLevel int) string {

	step := 4
	smb := " "
	indent := ""
	result := ""
	indent = strings.Repeat(smb, curLevel*step)
	for _, r := range diff {
		//1
		//Простые значений new/added
		if len(r.Value) == 1 {

			result += indent + getSymbol(r) + r.Key + ": " + getValue(r.Value[0], "stylish") + "\n"
			continue
		}

		//2.
		//updated простых и рекурсивных значений
		if len(r.Value) == 2 {

			//Для обновления вручную задаем знаки и отступы
			if nestedDiff, ok := r.Value[0].([]types.DiffItem); ok {
				//если первый элемент ! нода
				result += indent + "- " + r.Key + ": {\n" + formatStylish(nestedDiff, curLevel+1)
				result += indent + getSymbol(r) + "}\n"
			} else {
				//если первое значение - просто тип
				result += indent + "- " + r.Key + ": " + getValue(r.Value[0], "stylish") + "\n"
			}

			if nestedDiff, ok := r.Value[1].([]types.DiffItem); ok {
				//если первый элемент ! нода
				result += indent + "+ " + r.Key + ": {\n" + formatStylish(nestedDiff, curLevel+1)
				result += indent + getSymbol(r) + "}\n"
			} else {
				//если первое значение - просто тип
				result += indent + "+ " + r.Key + ": " + getValue(r.Value[1], "stylish") + "\n"
			}

			continue
		}
		//3 для вложенных элементов
		if len(r.Children) > 0 {
			//Для вложенных
			//fmt.Println(r.Key)

			indent = strings.Repeat(smb, curLevel*step)
			result += indent + getSymbol(r) + r.Key + ": {\n"
			result += formatStylish(r.Children, curLevel+1)
			result += indent + strings.Repeat(smb, 2) + "}\n"
		}

	}
	return result
}
