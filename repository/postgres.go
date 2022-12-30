package repository

import (
	"account-api/model"
	"gorm.io/gorm"
)

type Account interface {
	Save(account model.Account) (model.Account, error)
	Find(id uint) (model.Account, error)
	FindAll() ([]model.Account, error)
	Update(account model.Account) (model.Account, error)
	Delete(id uint) error
}

type User interface {
	Save(user model.User) (model.User, error)
	Find(id uint) (model.User, error)
	FindAll() ([]model.User, error)
	Delete(id uint) error
}

type AccountRepositoryImpl struct {
	db *gorm.DB
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) Account {
	return &AccountRepositoryImpl{db}
}

func NewUserRepository(db *gorm.DB) User {
	return &UserRepositoryImpl{db}
}

func (a *AccountRepositoryImpl) Save(account model.Account) (model.Account, error) {
	result := a.db.Create(&account)
	return account, result.Error
}

func (a *AccountRepositoryImpl) Find(id uint) (model.Account, error) {
	var account model.Account
	result := a.db.First(&account, id)
	return account, result.Error
}

func (a *AccountRepositoryImpl) FindAll() ([]model.Account, error) {
	var accounts []model.Account
	result := a.db.Find(&accounts)
	return accounts, result.Error
}

func (a *AccountRepositoryImpl) Update(account model.Account) (model.Account, error) {
	result := a.db.Save(&account)
	return account, result.Error
}

func (a *AccountRepositoryImpl) Delete(id uint) error {
	return a.db.Delete(&model.Account{}, id).Error
}

func (u *UserRepositoryImpl) Save(user model.User) (model.User, error) {
	result := u.db.Create(&user)
	return user, result.Error
}

func (u *UserRepositoryImpl) Find(id uint) (model.User, error) {
	var user model.User
	result := u.db.First(&user, id)
	return user, result.Error
}

func (u *UserRepositoryImpl) FindAll() ([]model.User, error) {
	var users []model.User
	result := u.db.Find(&users)
	return users, result.Error
}

func (u *UserRepositoryImpl) Delete(id uint) error {
	//TODO implement me
	panic("implement me")
}
