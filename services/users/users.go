package users

import (
	"context"
	"github.com/shijting/web/protos"
)

type UserServiceImpl struct {}

func (*UserServiceImpl) Register(ctx context.Context, req *protos.UserRegisterRequest) (*protos.UserRegisterResponse, error) {

	return &protos.UserRegisterResponse{Code: 200, Msg: req.GetUsername() + "---"+req.Email}, nil
}