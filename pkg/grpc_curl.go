package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"sort"
	"strings"
)

type GrpcCurl struct {
	addr string
	ctx  context.Context

	client *grpcreflect.Client
}

func NewGrpcCurl(ctx context.Context, addr string) (*GrpcCurl, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	refClient := grpcreflect.NewClientAuto(ctx, conn)
	return &GrpcCurl{
		addr:   addr,
		ctx:    ctx,
		client: refClient,
	}, nil
}

func (g *GrpcCurl) ListMethods(serviceName string) ([]string, error) {
	dsc, err := g.FindSymbol(serviceName)
	if err != nil {
		return nil, err
	}
	if sd, ok := dsc.(*desc.ServiceDescriptor); !ok {
		return nil, errors.New(fmt.Sprintf("not found service name %s", serviceName))
	} else {
		methods := make([]string, 0, len(sd.GetMethods()))
		for _, method := range sd.GetMethods() {
			methods = append(methods, method.GetFullyQualifiedName())
		}
		sort.Strings(methods)
		return methods, nil
	}
}

func (g *GrpcCurl) ListServices() ([]string, error) {
	services, err := g.client.ListServices()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, service := range services {
		if strings.HasPrefix(service, "grpc") {
			continue
		}
		if strings.HasPrefix(service, "kratos") {
			continue
		}
		result = append(result, service)
	}
	return result, nil
}

func (g *GrpcCurl) FindSymbol(fullyQualifiedName string) (desc.Descriptor, error) {
	file, err := g.client.FileContainingSymbol(fullyQualifiedName)
	if err != nil {
		return nil, err
	}
	d := file.FindSymbol(fullyQualifiedName)
	if d == nil {
		return nil, errors.New(fmt.Sprintf("could not find symbol: %s", fullyQualifiedName))
	}
	return d, nil
}
