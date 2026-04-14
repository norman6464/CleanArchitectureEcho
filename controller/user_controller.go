package controller

import "github.com/labstack/echo"

// Controllerインターフェースを作成してUsecaseの依存関係逆転させる（依存関係の逆転の原則）
type IUserController interface {
	Signup(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
}
