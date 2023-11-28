package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sort"
	"strings"
)

type GrpcCurl struct {
	addr string
	ctx  context.Context

	client *grpcreflect.Client

	methodMap map[string]*desc.MethodDescriptor
	conn      *grpc.ClientConn
}

func NewGrpcCurl(ctx context.Context, addr string) (*GrpcCurl, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	refClient := grpcreflect.NewClientAuto(ctx, conn)
	return &GrpcCurl{
		conn:      conn,
		addr:      addr,
		ctx:       ctx,
		client:    refClient,
		methodMap: make(map[string]*desc.MethodDescriptor),
	}, nil
}

func (g *GrpcCurl) GetConn() *grpc.ClientConn {
	return g.conn
}

func (g *GrpcCurl) ListMethods(serviceName string) ([]string, error) {
	dsc, err := g.FindSymbol(serviceName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sd, ok := dsc.(*desc.ServiceDescriptor)
	if !ok {
		return nil, errors.New(fmt.Sprintf("not found service name %s", serviceName))
	}
	result := make([]string, 0, len(sd.GetMethods()))
	for _, method := range sd.GetMethods() {
		methodName := method.GetFullyQualifiedName()
		g.methodMap[methodName] = method
		result = append(result, methodName)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i] > result[j]
	})
	return result, nil
}

func (g *GrpcCurl) GetMethodDescByName(methodName string) (*desc.MethodDescriptor, error) {
	if len(g.methodMap) == 0 { // init current service methods¬
		arr := strings.Split(methodName, ".")
		_, err := g.ListMethods(strings.Join(arr[0:len(arr)-1], "."))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	r, ok := g.methodMap[methodName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("method not found: %s", methodName))
	}
	return r, nil
}

func (g *GrpcCurl) genPayload(method *desc.MethodDescriptor) string {
	payload := make(map[string]any)
	for _, field := range method.GetInputType().GetFields() {
		payload[field.GetName()] = DefaultFieldValue(field)
	}
	bytes, err := json.Marshal(&payload)
	if err != nil {
		log.Errorf("json.Marshal error: %v", err)
	}
	return string(bytes)
}

func (g *GrpcCurl) ListServices() ([]string, error) {
	services, err := g.client.ListServices()
	if err != nil {
		return nil, errors.WithStack(err)
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
		return nil, errors.WithStack(err)
	}
	d := file.FindSymbol(fullyQualifiedName)
	if d == nil {
		return nil, errors.New(fmt.Sprintf("could not find symbol: %s", fullyQualifiedName))
	}
	return d, nil
}
