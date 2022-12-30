package service

import (
	"account-api/model"
	"account-api/repository"
)

type AccountService interface {
	CreateAccount(account model.Account) (model.Account, error)
	GetAccount(id uint) (model.Account, error)
	GetAccounts() ([]model.Account, error)
	UpdateAccount(account model.Account) (model.Account, error)
	DeleteAccount(id uint) error
}

type AccountServiceImpl struct {
	AccountRepository repository.Account
}

func NewAccountService(accountRepository repository.Account) AccountService {
	return &AccountServiceImpl{accountRepository}
}

func (a *AccountServiceImpl) CreateAccount(account model.Account) (model.Account, error) {
	return a.AccountRepository.Save(account)
}

func (a *AccountServiceImpl) GetAccount(id uint) (model.Account, error) {
	return a.AccountRepository.Find(id)
}

func (a *AccountServiceImpl) GetAccounts() ([]model.Account, error) {
	return a.AccountRepository.FindAll()
}

func (a *AccountServiceImpl) UpdateAccount(account model.Account) (model.Account, error) {
	return a.AccountRepository.Update(account)
}

func (a *AccountServiceImpl) DeleteAccount(id uint) error {
	return a.AccountRepository.Delete(id)
}
