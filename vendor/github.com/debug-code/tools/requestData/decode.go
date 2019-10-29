package requestData

import "strings"

func RequestBase64ToBase64(rb64 string) string {
	if rb64 == "" {
		return ""
	}
	data := strings.Split(string(rb64), ",")
	str := data[1] + "="
	return str
}
