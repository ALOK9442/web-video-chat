package helpers

import (
	"encoding/json"
)

func MarshalMessage(t string, data map[string]interface{}) []byte {
	outerMessage := map[string]interface{}{
		"type": t,
		"data": data,
	}
	finalJSON, _ := json.Marshal(outerMessage)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	return finalJSON
}
