package pkg

import (
	"context"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
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

func (g *GrpcCurl) ListMethods() ([]string, error) {
	return nil, nil
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
