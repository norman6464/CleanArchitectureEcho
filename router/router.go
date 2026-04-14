package router

import (
	"go-rest-api/controller"

	"github.com/labstack/echo/v4"
)

// DI（依存性注入）
func NewRouter(uc controller.IUserController) *echo.Echo {
	e := echo.New() // echoのインスタンス化をする
	e.POST("/signup", uc.Signup)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	return e
}
