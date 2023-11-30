package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

var gcurl *GrpcCurl
var ctx context.Context

func init() {
	ctx = context.TODO()
	var err error
	gcurl, err = NewGrpcCurl(ctx, "10.0.81.250:15555")
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
			t.Log(method)
		}
	}
}

func TestMethodPayload(t *testing.T) {
	methodName := "tophant.sensitive_info.api.SensitiveInfoService.List"
	method, err := gcurl.GetMethodDescByName(methodName)
	if err != nil {
		t.Fatal(err)
	}
	payload := gcurl.genPayload(method)

	// Marshal the map into a JSON-formatted string
	formattedJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	t.Logf("\n%s\n", formattedJSON)
}

func TestLoopUpHost(t *testing.T) {

}
