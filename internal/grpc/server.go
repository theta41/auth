package grpc

import (
	"context"
	"fmt"
	"net"

	"gitlab.com/g6834/team41/auth/internal/env"
	"google.golang.org/grpc"

	pb "gitlab.com/g6834/team41/auth/api/auth"
)

func StartServer(host string) error {

	lis, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("grpc failed to listen: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})

	env.E().L.Printf("grpc server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("grpc failed to serve: %w", err)
	}

	return nil
}

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateResponse, error) {

	env.E().L.Printf("grpc.Validate %v", in)

	//return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")

	return &pb.ValidateResponse{
		Success:      true,
		Login:        in.Login,
		AccessToken:  in.AccessToken,
		RefreshToken: in.RefreshToken,
	}, nil
}
