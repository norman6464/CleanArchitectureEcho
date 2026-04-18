package router

import (
	"go-rest-api/controller"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4" // 別名(エイリアス)

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// DI（依存性注入）
func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	e := echo.New() // echoのインスタンス化をする

	// ミドルウェアの作成
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))

	// Spring bootでは自動で付与されていなかったのでGoのEchoでは自分で設定をする必要がある
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		// CookieSameSite: http.SameSiteMode,
		CookieSameSite: http.SameSiteDefaultMode, // Postmanの使用時だけhttp.SameSiteDefaultModeにしておく
		// CookieMaxAge: 60
	}))

	e.POST("/signup", uc.Signup)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)
	t := e.Group("/tasks")
	t.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	t.GET("", tc.GetAllTasks)
	t.GET("/:taskId", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:taskId", tc.UpdateTask)
	t.DELETE("/:taskId", tc.DeleteTask)
	return e
}
