package helpers

import (
	"encoding/json"
	"io"
	"strings"
)

func TrimSpaceToUpper(input string) string {
	return strings.TrimSpace(strings.ToUpper(input))
}

func PrettyJSON(input interface{}) string {
	bytes, _ := json.MarshalIndent(input, "", "    ")
	return string(bytes)
}
func WriteToWriter(w io.Writer, input interface{}) {
	w.Write(ToJSONBytes(input))
}
func ToJSONBytes(input interface{}) []byte {
	jsonBytes, _ := json.Marshal(input)
	return jsonBytes
}
