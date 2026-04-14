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

// Spring bootの時はアノテーションやコンストラクタインジェクションで意識をしていなかったのだがインターフェースを返却をするが実際には実装をしたrepositoryが帰っている
// これが依存性逆転の原則になる実際にこのコンストラクタを使用することでUsecaseが依存するのがこのインターフェースのリポジトリになる
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
