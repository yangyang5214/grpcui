package pkg

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type GrpcSend struct {
	conn *grpc.ClientConn
}

func NewGrpcSend(conn *grpc.ClientConn) *GrpcSend {
	return &GrpcSend{
		conn: conn,
	}
}

func (g *GrpcSend) send(ctx context.Context, method *desc.MethodDescriptor, request proto.Message) (string, error) {
	stub := grpcdynamic.NewStubWithMessageFactory(g.conn, nil)
	resp, err := stub.InvokeRpc(ctx, method, request)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return resp.String(), nil
}
