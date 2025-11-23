package formatters

import (
	"code/types"
	"strings"
)

func formatPlain(diff []types.DiffItem) string {

	result := getFormatArray(diff)
	return strings.Join(result, "\n")
}

func getFormatArray(diff []types.DiffItem) []string {
	result := []string{}
	value1 := ""
	value2 := ""
	for _, r := range diff {

		switch r.Result {
		case "deleted":
			result = append(result, "Property '"+r.Path+"' was removed")

		case "new":

			if len(r.Children) > 0 {
				value1 = "[complex value]"
			} else {
				value1 = getValue(r.Value[0], "plain")
			}
			result = append(result, "Property '"+r.Path+"' was added with value: "+value1)

		case "updated":

			if _, ok := r.Value[0].([]types.DiffItem); ok {
				value1 = "[complex value]"
			} else {
				value1 = getValue(r.Value[0], "plain")
			}

			if _, ok := r.Value[1].([]types.DiffItem); ok {
				value2 = "[complex value]"
			} else {
				value2 = getValue(r.Value[1], "plain")

			}

			result = append(result, "Property '"+r.Path+"' was updated. From "+value1+" to "+value2)

		default:
			//result = append(result, []string{})
		}

		if len(r.Children) > 0 {
			result = append(result, getFormatArray(r.Children)...)
		}

	}

	return result
}
