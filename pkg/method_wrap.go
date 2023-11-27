package pkg

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/descriptorpb"
	"reflect"
)

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

func JsonToMessage(jsonData map[string]any, message *dynamic.Message) error {
	var err error
	for key, value := range jsonData {
		fd := message.FindFieldDescriptorByName(key)

		finalValue := TransFormField(fd, reflect.ValueOf(value))
		if finalValue == nil {
			return errors.New(fmt.Sprintf("field %s transforme failed, %v", key, fd.GetMessageType()))
		}
		err = message.TrySetField(fd, finalValue)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func TransFormField(fd *desc.FieldDescriptor, val reflect.Value) any {
	t := fd.GetType()
	switch t {
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return int32(val.Int())

	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return val.Int()

	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		return uint32(val.Uint())

	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		return val.Uint()

	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return float32(val.Float())

	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return val.Float()

	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return val.Bool()

	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return val.Bytes()

	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return val.String()

	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		//todo
		return nil
	default:
		return nil
	}
}
