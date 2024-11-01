package interfaces

import (
	pb "github.com/shivaraj-shanthaiah/user-management/pkg/proto"
)

type UserServiceInter interface {
	CreateUserService(p *pb.Create) (*pb.Response, error)
	GetUserByIDService(p *pb.ID) (*pb.Profile, error)
	UpdateUserService(profile *pb.Profile) (*pb.Profile, error)
	DeleteuserByIDService(Id *pb.ID) (*pb.Response, error)
	GetAllUsersService(p *pb.NoParams) (*pb.Names, error)
}
