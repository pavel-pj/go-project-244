package formatters

import "code/types"

func formatPlain(diff []types.DiffItem) string {
	result := ""
	value1 := ""
	value2 := ""
	for _, r := range diff {

		switch r.Result {
		case "deleted":
			result += "Property '" + r.Path + "' was removed\n"
		case "new":

			if len(r.Children) > 0 {
				value1 = "[complex value]"
			} else {
				value1 = getValue(r.Value[0], "plain")
			}
			result += "Property '" + r.Path + "' was added with value: " + value1 + "\n"

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

			result += "Property '" + r.Path + "' was updated. From " + value1 + " to " + value2 + "\n"

		default:
			result += ""
		}

		if len(r.Children) > 0 {
			result += formatPlain(r.Children)
		}

	}
	return result
}
