package pkg

import (
	"reflect"
)

func anyToInt32(v any) int32 {
	val := reflect.ValueOf(v)
	if val.CanInt() {
		return int32(val.Int())
	}
	return 0 //default
}
