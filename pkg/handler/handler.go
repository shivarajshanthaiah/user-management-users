package handler

import (
	"context"

	pb "github.com/shivaraj-shanthaiah/user-management/pkg/proto"
	inter "github.com/shivaraj-shanthaiah/user-management/pkg/service/interfaces"
)

type UserHandler struct {
	SVC inter.UserServiceInter
	pb.UserServiceServer
}

func NewUserHandler(svc inter.UserServiceInter) *UserHandler {
	return &UserHandler{
		SVC: svc,
	}
}

func (u *UserHandler) CreateUser(ctx context.Context, p *pb.Create) (*pb.Response, error) {
	response, err := u.SVC.CreateUserService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (u *UserHandler) GetUserByID(ctx context.Context, p *pb.ID) (*pb.Profile, error) {
	response, err := u.SVC.GetUserByIDService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (u *UserHandler) UpdateUser(ctx context.Context, p *pb.Profile) (*pb.Profile, error) {
	response, err := u.SVC.UpdateUserService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (u *UserHandler) DeleteUserBYID(ctx context.Context, p *pb.ID) (*pb.Response, error) {
	response, err := u.SVC.DeleteuserByIDService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (u *UserHandler) GetAllUsers(ctx context.Context, p *pb.NoParams) (*pb.Names, error) {
	response, err := u.SVC.GetAllUsersService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
