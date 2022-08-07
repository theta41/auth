package grpc

import (
	"context"
	"fmt"
	"net"

	"gitlab.com/g6834/team41/auth/internal/models"
	"gitlab.com/g6834/team41/auth/internal/ports"

	"gitlab.com/g6834/team41/auth/internal/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "gitlab.com/g6834/team41/auth/api/auth"
)

func StartServer(host string, auth ports.Auth) error {

	lis, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("grpc failed to listen: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{auth: auth})

	env.E().L.Printf("grpc server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("grpc failed to serve: %w", err)
	}

	return nil
}

type server struct {
	auth ports.Auth
	pb.UnimplementedAuthServiceServer
}

func (s *server) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateResponse, error) {

	env.E().L.Printf("grpc.Validate %v", in)

	//return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")

	oldTokens := models.TokenPair{
		AccessToken:  in.AccessToken,
		RefreshToken: in.RefreshToken,
	}
	newTokens, err := s.auth.Validate(in.Login, oldTokens)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Validate %v", err)
	}

	return &pb.ValidateResponse{
		Success:      true,
		Login:        in.Login,
		AccessToken:  newTokens.AccessToken,
		RefreshToken: newTokens.RefreshToken,
	}, nil
}
