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
		methods, mErr := gcurl.ListMethods(service)
		if mErr != nil {
			t.Fatal(mErr)
		}
		t.Log(methods)
	}
}
