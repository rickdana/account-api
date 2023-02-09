package service

import (
	"account-api/model"
	"account-api/repository"
	"context"
	"github.com/google/uuid"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUser(id uint) (model.User, error)
	GetUsers() ([]model.User, error)
	UpdateUser(user model.User) (model.User, error)
	DeleteUser(id uint) error
}

type UserServiceImpl struct {
	UserRepository repository.User
	eventSender    EventSender
}

func NewUserServiceImpl(eventSender EventSender, userRepository repository.User) UserService {
	return &UserServiceImpl{eventSender: eventSender, UserRepository: userRepository}
}

func (u *UserServiceImpl) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	user.UpdateBy = "system"
	user.ID = uint(uuid.New().ID())
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	messageKey := model.NewFourEyesMessageKey(user.ID, "user", "create")
	msg := model.NewFourEyesMessage(nil, user, "system")

	err := u.eventSender.Send(*messageKey, *msg)
	if err != nil {
		return model.User{}, err
	}
	return model.User{}, nil
}

func (u *UserServiceImpl) GetUser(id uint) (model.User, error) {
	return u.UserRepository.Find(id)
}

func (u *UserServiceImpl) GetUsers() ([]model.User, error) {
	return u.UserRepository.FindAll()
}

func (u *UserServiceImpl) UpdateUser(user model.User) (model.User, error) {
	return u.UserRepository.Save(user)
}

func (u *UserServiceImpl) DeleteUser(id uint) error {
	return u.UserRepository.Delete(id)
}
