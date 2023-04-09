package util

import "strings"

// CheckString 检查string
func CheckString(stringValue string) string {
	// 去除空格
	stringValue = strings.Replace(stringValue, " ", "", -1)
	// 去除换行符
	return strings.Replace(stringValue, "\n", "", -1)
}
