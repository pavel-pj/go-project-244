package formatters

import (
	"code/types"
	"encoding/json"
)

func formatJson(diff []types.DiffItem) string {

	// Convert struct to JSON bytes
	jsonData, err := json.Marshal(diff)
	if err != nil {
		panic(err)
	}
	return string(jsonData)

}
