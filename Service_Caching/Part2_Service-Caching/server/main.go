package main

import (
	"context"
	"errors"
	"net"
	proto "part2/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterUserServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}

}
func (s *server) GetUserByName(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	// name := request.GetName()
	// stlen := strconv.Itoa(len(name))
	// return &proto.Response{Name: name, Roll: int64(len(name) * 10), Class: stlen}, nil
	return nil, errors.New("not implemented yet. Rupin will implement me")
}
