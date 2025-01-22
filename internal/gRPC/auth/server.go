package auth

import (
	"context"
	"errors"
	"gRPCAuthService/internal/services/auth"
	"gRPCAuthService/internal/storage"
	ssov1 "github.com/21Timofei/Contracts/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}
type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

const (
	emptyValue = 0
)

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil

}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateAdminState(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLoginRequest(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}

	if req.AppId == emptyValue {
		return status.Error(codes.InvalidArgument, "app required")
	}

	return nil
}

func validateRegisterRequest(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}

	return nil
}

func validateAdminState(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user required")
	}
	return nil
}
