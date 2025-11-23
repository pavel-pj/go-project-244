package formatters

import (
	"code/types"
	"fmt"
	"strconv"
)

func Formater(diff []types.DiffItem, format string) string {

	level := 0
	switch format {
	case "plain":
		return formatPlain(diff)
	case "stylish":
		return "{\n" + formatStylish(diff, level) + "}"
	case "json":
		return formatJson(diff)
	default:
		return "{\n" + formatStylish(diff, level) + "}"

	}

}

func getSymbol(item types.DiffItem) string {
	switch item.Result {
	case "new":
		return "+ "
	case "deleted":
		return "- "
	case "unchanged":
		return "  "
	}
	return "  "
}

func getValue(value interface{}, format string) string {
	switch v := value.(type) {
	case string:
		if format == "plain" {
			return "'" + v + "'"
		}
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
