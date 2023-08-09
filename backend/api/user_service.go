package api

import (
	connect "connectrpc.com/connect"
	"context"
	v1 "github.com/filebrowser/filebrowser/gen/proto/filebrowser/v1"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u UserService) List(ctx context.Context, c *connect.Request[v1.UserServiceListRequest]) (*connect.Response[v1.UserServiceListResponse], error) {
	//TODO implement me
	panic("implement me")
}
