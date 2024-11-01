package repo

import (
	"errors"

	"github.com/shivaraj-shanthaiah/user-management/pkg/model"
	inter "github.com/shivaraj-shanthaiah/user-management/pkg/repo/interfaces"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) inter.UserRepoInter {
	return &UserRepo{
		DB: db,
	}
}

func (u *UserRepo) CreateUser(user *model.User) (uint, error) {
	if err := u.DB.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (u *UserRepo) FindUserByID(userID uint32) (*model.User, error) {
	var user model.User

	if err := u.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) UpdateUser(user *model.User) error {
	if err := u.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) DeleteUserByID(userID uint) error {
	if err := u.DB.Delete(&model.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) FindUserByEmail(email string) (*model.User, error) {
	var user model.User

	if err := u.DB.Model(&model.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) FindUserByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := u.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) FindAllUsers() (*[]model.User, error) {
	var users []model.User

	if err := u.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}
