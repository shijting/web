package users

import (
	"context"
	"github.com/jinzhu/copier"
	dao "github.com/shijting/web/dao/users"
	"github.com/shijting/web/models/users"
	"github.com/shijting/web/protos"
	"github.com/shijting/web/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type UserServiceImpl struct {
	mu    *sync.RWMutex
}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{mu: &sync.RWMutex{}}
}

// 用户注册rpc服务实现
func (u *UserServiceImpl) Register(ctx context.Context, req *protos.UserRegisterRequest) (*protos.UserResponse, error) {
	//u.mu.Lock()
	//defer u.mu.Unlock()
	//if err :=req.Validate();err !=nil {
	//	return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	//}
	if req.GetPassword() != req.GetPasswordRepeat() {
		return nil, status.Errorf(codes.InvalidArgument, "两次密码不一致")
	}
	passwd, err := services.HashPassword(req.GetPassword())
	if err !=nil {
		return nil, status.Errorf(codes.InvalidArgument, "register failed")
	}
	data := &users.User{
		Username:    req.GetUsername(),
		Email:       req.GetEmail(),
		Password:    passwd,
		CreatedDate: time.Now().UTC(),
	}
	result,err := dao.Insert(nil, data)
	if err !=nil {
		return nil, status.Errorf(codes.InvalidArgument, "register failed")
	}
	userInfo := &protos.User{}
	copier.Copy(userInfo, result)
	return &protos.UserResponse{Code: int32(codes.OK), Message: "success", Details: userInfo}, nil
}

// 用户登录
func (u *UserServiceImpl) Login(ctx context.Context, req *protos.UserLoginRequest) (*protos.UserResponse, error) {
	//if err :=req.Validate();err !=nil {
	//	return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	//}
	result,err := dao.GetOneByUsername(nil, req.GetUsername())
	if err !=nil {
		return nil, status.Errorf(codes.InvalidArgument, "登录失败: 用户名或密码不正确")
	}
	if !services.CheckPasswordHash(req.GetPassword(), result.Password){
		return nil, status.Errorf(codes.InvalidArgument, "登录失败: 用户名或密码不正确")
	}
	// 登录成功
	// TODO generate token
	// ......
	userInfo := &protos.User{}
	copier.Copy(userInfo, result)
	return &protos.UserResponse{Code: int32(codes.OK), Message: "登录成功", Details: userInfo}, nil
}