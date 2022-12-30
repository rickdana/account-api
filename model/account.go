package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	UpdateBy string `json:"update_by"`
}

func (u *User) Update(user User) {
	u.Username = user.Username
	u.Password = user.Password
	u.Email = user.Email
	u.Role = user.Role
}

type Account struct {
	gorm.Model
	AccountNumber  string  `json:"account_number" validate:"required" gorm:"uniqueIndex"`
	AccountName    string  `json:"account_name" validate:"required"`
	AccountType    string  `json:"account_type" validate:"required"`
	AccountBalance float64 `json:"account_balance" validate:"required"`
	Currency       string  `json:"currency" validate:"required"`
	UpdatedBy      string  `json:"updated_by"`
}

func NewAccount() *Account {
	a := Account{}
	a.ID = uint(uuid.New().ID())
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	return &a
}

func (a *Account) Update(account Account) {
	a.AccountNumber = account.AccountNumber
	a.AccountName = account.AccountName
	a.AccountType = account.AccountType
	a.AccountBalance = account.AccountBalance
	a.Currency = account.Currency
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uint(uuid.New().ID())
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	u.UpdatedAt = now
	return nil
}
