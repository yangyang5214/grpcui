package pkg

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/descriptorpb"
	"reflect"
)

type MethodWrap struct {
	Service string   `json:"service,omitempty"`
	Methods []string `json:"methods,omitempty"`
}

type AllMethod struct {
	Addr    string        `json:"addr,omitempty"`
	Methods []*MethodWrap `json:"methods,omitempty"`
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
		if fd == nil {
			return errors.New(fmt.Sprintf("key %s is nil", key))
		}

		finalValue := TransFormField(fd, key, value)
		if finalValue == nil {
			return errors.New(fmt.Sprintf("field %s transforme failed, %v", key, fd.GetMessageType()))
		}
		log.Infof("transform <%v> value from <%v> to <%v>", key, value, finalValue)
		err = message.TrySetField(fd, finalValue)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func TransFormField(fd *desc.FieldDescriptor, key string, val any) any {
	if val == nil {
		return nil
	}
	v := reflect.ValueOf(val)
	t := fd.GetType()
	switch t {
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return int32(v.Int())

	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return v.Int()

	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		return uint32(v.Uint())

	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		return v.Uint()

	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return float32(v.Float())

	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return v.Float()

	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return v.Bool()

	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return v.Bytes()

	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return v.String()

	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		//todo
		return nil
	default:
		return nil
	}
}
