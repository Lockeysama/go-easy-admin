package utils

import "reflect"

// Contain 查询 value 是否在 in 切片中
func Contain(value interface{}, in *[]interface{}) bool {
	valueRef := reflect.TypeOf(value)
	for _, inSetValue := range *in {
		inSetValueRef := reflect.TypeOf(inSetValue)
		if valueRef.Kind() != inSetValueRef.Kind() {
			return false
		}
		if value == inSetValue {
			return true
		}
	}
	return false
}

// StringContain 查询 value 是否在 in 切片中
func StringContain(value string, in *[]string) bool {
	for _, inSetValue := range *in {
		if value == inSetValue {
			return true
		}
	}
	return false
}
