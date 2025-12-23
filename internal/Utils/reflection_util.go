package utils

import (
	"reflect"
	"strings"
)

var cache = make(map[string][]string)

// given a value of struct returns list of jsontags of fields
func GetFieldsJSONTags(u any) []string {
	t := reflect.TypeOf(u)
	if tags, exist := cache[t.Name()]; exist {
		return tags
	} else {
		n := t.NumField()
		fields := make([]string, 1)
		for i := 0; i < n; i++ {
			jsonTag := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
			if len(jsonTag) != 0 {
				fields = append(fields, jsonTag)
			}
		}
		cache[t.Name()] = fields
		return fields
	}
}
