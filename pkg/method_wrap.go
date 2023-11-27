package pkg

import "github.com/jhump/protoreflect/desc"

type MethodWrap struct {
	method  string
	payload string //json payload
}

func DefaultFieldValue(field *desc.FieldDescriptor) any {
	name := field.GetName()
	switch name {
	case "page":
		return 1
	case "page_size", "size":
		return 10
	}
	return field.GetDefaultValue()
}
