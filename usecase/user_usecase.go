package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	Signup(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error) // string型の戻り値はJWTトークンの値を返す
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

// コンストラクタインジェクション
func NewuserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

// サインアップメソッド
func (uu *userUsecase) Signup(user model.User) (model.UserResponse, error) {

	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err // {}はインスタンス化であり、そこでまずは0値で返却をしている
	}
	newUser := model.User{Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err // これも0値で返しつつ、errも返している
	}

	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}

	return resUser, nil
}

// ログインメソッド
func (uu *userUsecase) Login(user model.User) (string, error) {

	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}

	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil { // ここでポインタを渡しているので内部的にstoredUserインスタンスは値が入っている状態になっている
		return "", err
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, err

}
