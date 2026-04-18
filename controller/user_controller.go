package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

// Controllerインターフェースを作成してUsecaseの依存関係逆転させる（依存関係の逆転の原則）
type IUserController interface {
	Signup(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
}

// Controller（interface層）からusecaseInterfaceに依存をしていく
// 内側に依存をしているがそもそも保守容易生保つためにusecaseもインターフェースの定義をしている
type userController struct {
	uu usecase.IUserUsecase
}

// DI（依存性注入）をしていき、コントローラーインターフェースが戻り値だが実際の中身にはControllerクラスの実装したクラスが返却されている
// 引数には依存先（今回の場合はコントローラーがusecaseインターフェースに依存をするのでIUsecase）をよび出している
func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu} // インスタンス化をしておりこれを呼び出してまずはインスタンス化をしている状態
}

// userControllerのアンドで渡しているのでその中身で代入をしたアドレス演算子でこちらの構造体に反映される
func (uc *userController) Signup(c echo.Context) error {
	// リクエストが正当なのかの確認をしていく
	user := model.User{} // インスタンス化
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	// サインアップメソッド
	userRes, err := uc.uu.Signup(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{} // インスタンス化
	// リクエストが構造体userにちゃんと当てはまっているかの検証
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	tokenString, err := uc.uu.Login(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	// サーバーサイドでJWTトークンをCookieに保存をすること
	cookie := new(http.Cookie) // Cookieインスタンス生成をしていく
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/" // 全てのパスでcookieを有効にしていく
	cookie.Domain = os.Getenv("API_DOMAIN")

	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode // フロントエンド、バックエンドともに違うドメインなのでSameSiteNoneModeにする
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK) // 何もボディに返却がない場合はNocontentメソッドを使用する
}

func (us *userController) LogOut(c echo.Context) error {

	// サーバーサイドでJWTトークンをCookieに保存をすること
	cookie := new(http.Cookie) // Cookieインスタンス生成をしていく
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")

	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode // フロントエンド、バックエンドともに違うドメインなのでSameSiteNoneModeにする
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK) // 何もボディに返却がない場合はNocontentメソッドを使用する
}

func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}
