package pkg

import (
	"context"
	"testing"
)

var gcurl *GrpcCurl
var ctx context.Context

func init() {
	ctx = context.TODO()
	var err error
	gcurl, err = NewGrpcCurl(ctx, "10.0.81.250:9000")
	if err != nil {
		panic(err)
	}
}

func TestListMethods(t *testing.T) {
	services, err := gcurl.ListServices()
	if err != nil {
		t.Fatal(err)
	}
	for _, service := range services {
		t.Logf("service is %s", service)
		methods, mErr := gcurl.ListMethods(service)
		if mErr != nil {
			t.Fatal(mErr)
		}
		for _, method := range methods {
			t.Log(method.method)
		}
	}
}

func TestMethodPayload(t *testing.T) {
	methodName := "tophant.parser.api.TaskResultService.List"
	method, err := gcurl.GetMethodDescByName(methodName)
	if err != nil {
		t.Fatal(err)
	}
	payload := gcurl.genPayload(method)
	t.Log(payload)
}
