package rpc

import (
	"context"

	userPb "github.com/filebrowser/filebrowser/v3/gen/proto/user/v1"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u UserService) Find(ctx context.Context, request *userPb.FindRequest) (*userPb.FindResponse, error) {
	panic("implement me")
}
