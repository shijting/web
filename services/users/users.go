package users

import (
	"context"
	"github.com/shijting/web/protos"
)

type UsersServiceImpl struct {}

func (*UsersServiceImpl) Register(ctx context.Context, req *protos.UserRegisterRequest) (*protos.UserRegisterResponse, error) {

	return nil, nil
}