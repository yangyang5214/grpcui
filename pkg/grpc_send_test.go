package pkg

import (
	"github.com/jhump/protoreflect/dynamic"
	"testing"
)

func TestSend(t *testing.T) {
	gSend := NewGrpcSend(gcurl.GetConn())
	methodName := "tophant.parser.api.TaskResultService.List"
	method, err := gcurl.GetMethodDescByName(methodName)
	if err != nil {
		t.Fatal(err)
	}

	jsonData := map[string]any{
		"page":    1,
		"size":    10,
		"task_id": "ASPM1729027707203571712",
	}

	dMsg := dynamic.NewMessage(method.GetInputType())
	err = JsonToMessage(jsonData, dMsg)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := gSend.send(ctx, method, dMsg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}
