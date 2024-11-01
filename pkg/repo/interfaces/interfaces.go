package interfaces

import "github.com/shivaraj-shanthaiah/user-management/pkg/model"

type UserRepoInter interface {
	CreateUser(user *model.User) (uint, error)
	FindUserByID(userID uint32) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUserByID(userID uint) error
	FindAllUsers() (*[]model.User, error)
	FindUserByEmail(email string) (*model.User, error)
	FindUserByPhone(phone string) (*model.User, error)
}
