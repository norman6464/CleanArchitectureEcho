package repository

import (
	"go-rest-api/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

// これはrepositoryがDB側に依存をしているのでDI（依存性の注入）を行なっている
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

// ur引数はメンバにDB変数を抱えているのでgormライブラリを使用できる
func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {

	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
