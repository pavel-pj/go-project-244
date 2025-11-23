package formatters

import (
	"code/types"
	"encoding/json"
)

func formatJson(diff []types.DiffItem) string {

	result := struct {
		Diff []types.DiffItem `json:"diff"`
	}{
		Diff: diff,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	return string(jsonData)

}
