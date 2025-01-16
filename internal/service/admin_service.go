package service

import (
	"context"

	"github.com/mohamedfawas/rmshop-admin-service/internal/domain"
	adminv1 "github.com/mohamedfawas/rmshop-proto/gen/v1/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type adminService struct {
	repo domain.AdminRepository
	adminv1.UnimplementedAdminServiceServer
}

func NewAdminService(repo domain.AdminRepository) adminv1.AdminServiceServer {
	return &adminService{repo: repo}
}

func (s *adminService) GetUserDetails(ctx context.Context, req *adminv1.GetUserDetailsRequest) (*adminv1.GetUserDetailsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	user, err := s.repo.GetUserDetails(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user details: %v", err)
	}

	return &adminv1.GetUserDetailsResponse{
		User: &adminv1.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}
