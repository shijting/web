package users

import "context"

type usersServiceImpl struct {}

func (*usersServiceImpl) Register(ctx context.Context, req *UserRegisterRequest) (*UserRegisterResponse, error) {

	return nil, nil
}