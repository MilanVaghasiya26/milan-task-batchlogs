package helpers

import (
	"reflect"
)

func GetStructName(entity interface{}) string {
	if t := reflect.TypeOf(entity); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
