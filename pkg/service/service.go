package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shivaraj-shanthaiah/user-management/config"
	"github.com/shivaraj-shanthaiah/user-management/pkg/model"
	pb "github.com/shivaraj-shanthaiah/user-management/pkg/proto"
	inter "github.com/shivaraj-shanthaiah/user-management/pkg/repo/interfaces"
	"github.com/shivaraj-shanthaiah/user-management/pkg/service/interfaces"
)

type UserService struct {
	Repo  inter.UserRepoInter
	redis *config.RedisService
}

func NewUserService(repo inter.UserRepoInter, redis *config.RedisService) interfaces.UserServiceInter {
	return &UserService{
		Repo:  repo,
		redis: redis,
	}
}

func (u *UserService) CreateUserService(p *pb.Create) (*pb.Response, error) {

	user := model.User{
		UserName: p.User_Name,
		Email:    p.Email,
		Phone:    p.Phone,
	}

	existingUserEmail, err := u.Repo.FindUserByEmail(user.Email)
	if err != nil {
		return &pb.Response{
			Status:  pb.Response_ERROR,
			Message: "error checking user email",
			Payload: &pb.Response_Error{
				Error: err.Error(),
			},
		}, err
	}

	if existingUserEmail != nil {
		return &pb.Response{
			Status:  pb.Response_ERROR,
			Message: "email already exists",
		}, nil
	}

	existingUserPhone, err := u.Repo.FindUserByPhone(user.Phone)
	if err != nil {
		fmt.Println("error checking user phone:", err)
		return &pb.Response{
			Status:  pb.Response_ERROR,
			Message: "error checking user phone",
			Payload: &pb.Response_Error{
				Error: err.Error(),
			},
		}, err
	}

	if existingUserPhone != nil {
		return &pb.Response{
			Status:  pb.Response_ERROR,
			Message: "phone number already exists",
		}, nil
	}

	userId, err := u.Repo.CreateUser(&user)
	if err != nil {
		return &pb.Response{
			Status:  pb.Response_ERROR,
			Message: "error creating user in database",
			Payload: &pb.Response_Error{
				Error: err.Error(),
			},
		}, errors.New("unable to create user")
	}

	log.Printf("User created with ID: %v", userId)

	cacheKey := "user_" + strconv.Itoa(int(userId))

	userData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling user data for caching: %v", err)
	} else {
		err = u.redis.SetDataInRedis(cacheKey, userData, time.Minute*2)
		if err != nil {
			log.Printf("Error setting user data in Redis: %v", err)
		}
	}
	return &pb.Response{
		Status:  pb.Response_OK,
		Message: "User created successfully",
		Payload: &pb.Response_Data{
			Data: strconv.Itoa(int(userId)),
		},
	}, nil
}

func (u *UserService) GetUserByIDService(p *pb.ID) (*pb.Profile, error) {

	cacheKey := "user_" + strconv.Itoa(int(p.ID))
	cachedUser, err := u.redis.GetFromRedis(cacheKey)
	if err == nil {
		// Unmarshal cached data into a User struct if cache hit
		var user model.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			log.Printf("Error unmarshalling cached user data: %v", err)
		} else {
			return &pb.Profile{
				User_ID:   uint32(user.ID),
				User_Name: user.UserName,
				Email:     user.Email,
				Phone:     user.Phone,
				Created:   user.CreatedAt.Format(time.RFC3339),
				Updated:   user.UpdatedAt.Format(time.RFC3339),
				Message:   "Fetched from redis",
			}, nil
		}
	} else if err != redis.Nil {
		// Handle Redis error (other than missing key)
		log.Printf("Error accessing Redis: %v", err)
		return nil, err
	}

	// Cache miss: fetch user from the database
	user, err := u.Repo.FindUserByID(p.ID)
	if err != nil {
		return nil, err
	}

	// Cache the retrieved user data for future requests
	userData, err := json.Marshal(user)
	if err == nil {
		_ = u.redis.SetDataInRedis(cacheKey, userData, time.Minute*2)
	} else {
		log.Printf("Error marshalling user data for caching: %v", err)
	}

	return &pb.Profile{
		User_ID:   uint32(user.ID),
		User_Name: user.UserName,
		Email:     user.Email,
		Phone:     user.Phone,
		Created:   user.CreatedAt.Format(time.RFC3339),
		Updated:   user.UpdatedAt.Format(time.RFC3339),
		Message:   "Fetched from postgres",
	}, nil
}

func (u *UserService) UpdateUserService(profile *pb.Profile) (*pb.Profile, error) {
	user, err := u.Repo.FindUserByID(profile.User_ID)
	if err != nil {
		return nil, err
	}

	user.UserName = profile.User_Name
	user.Phone = profile.Phone
	if err := u.Repo.UpdateUser(user); err != nil {
		return nil, err
	}

	cacheKey := "user_" + strconv.Itoa(int(user.ID))
	userData, err := json.Marshal(user)
	if err == nil {
		_ = u.redis.SetDataInRedis(cacheKey, userData, time.Minute*2)
	} else {
		log.Printf("Error marshalling user data for caching: %v", err)
	}

	return &pb.Profile{
		User_ID:   uint32(user.ID),
		User_Name: user.UserName,
		Email:     user.Email,
		Phone:     user.Phone,
		Created:   user.CreatedAt.Format(time.RFC3339),
		Updated:   user.UpdatedAt.Format(time.RFC3339),
		Message:   "Updated in postgres and cached in redis",
	}, nil
}

func (u *UserService) DeleteuserByIDService(Id *pb.ID) (*pb.Response, error) {
	err := u.Repo.DeleteUserByID(uint(Id.ID))
	if err != nil {
		log.Printf("Error deleting user with ID %d: %v", Id.ID, err)
		return nil, err
	}

	_, err = u.redis.DelFromRedis("user_"+strconv.Itoa(int(Id.ID)), "users")
	if err != nil {
		log.Printf("Error deleting user from Redis with key user_%d: %v", Id.ID, err)
	}

	return &pb.Response{
		Status:  pb.Response_OK,
		Message: "User deleted successfully",
		Payload: &pb.Response_Data{
			Data: "",
		},
	}, nil
}

func (u *UserService) GetAllUsersService(p *pb.NoParams) (*pb.Names, error) {
	users, err := u.Repo.FindAllUsers()
	if err != nil {
		return nil, err
	}

	var userList pb.Names
	for _, user := range *users {
		userList.Users = append(userList.Users, &pb.Profile{
			User_ID:   uint32(user.ID),
			User_Name: user.UserName,
		})
	}
	return &userList, nil
}