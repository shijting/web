package users

import (
	"context"
	"github.com/shijting/web/protos"
	"sync"
)

type UserServiceImpl struct {
	mu    *sync.RWMutex
}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{mu: &sync.RWMutex{}}
}

type Ret struct {
	Id string `json:"id"`
	Age int `json:"age"`
}
func (*UserServiceImpl) Register(ctx context.Context, req *protos.UserRegisterRequest) (*protos.UserRegisterResponse, error) {
	ret := map[string]interface{}{"id": "123", "age": "12"}
	return &protos.UserRegisterResponse{Message: req.GetUsername() + "---"+req.GetEmail(), Data: ret}, nil
}


