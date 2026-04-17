package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Controllerはechoなどのフレームワーク（HTTPサーバーなど）を使用をする
type ITaskController interface {
	GetAllTasks(c echo.Context) error
	GetTaskById(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
}

// controllerの依存先はusecaseになる
type taskController struct {
	tu usecase.ITaskUsecase
}

func NewTaskController(tu usecase.ITaskUsecase) ITaskController {
	return &taskController{tu}
}

func (tc *taskController) GetAllTasks(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)    // 二番目のメソッドではGETメソッドの戻りのinterfaceからの型アサーションをしている
	claims := user.Claims.(jwt.MapClaims) // jwt.Claimsインターンフェースから型アサーションをしえしている
	userId := claims["user_id"]           // 型アサーションをしないuser.Claims["user_id"]これはエラーになる

	taskRes, err := tc.tu.GetAllTasks(uint(userId.(float64))) // JSONの数値は全てfloat64になるのでカタアサーションをしないといけない
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, taskRes)
}

func (tc *taskController) GetTaskById(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)    // interfaceを型アサーションに変換をしていく
	claims := user.Claims.(jwt.MapClaims) // TokenインターフェースをMapClaimsの型にアサーション
	userId := claims["user_id"]
	id := c.Param("taskId")       // パスパラメーターで取得をしている
	taskId, _ := strconv.Atoi(id) //　stringがたをint型にしている
	taskRes, err := tc.tu.GetTaskById(uint(userId.(float64)), uint(taskId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, taskRes)

}

func (tc *taskController) CreateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	task := model.Task{} // インスタンス化を先にしておく

	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// userIdは代入されていないので変数に入れる
	task.UserId = uint(userId.(float64)) // json形式ではfloat64にする
	taskRes, err := tc.tu.CreateTask(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, taskRes)

}

func (tc *taskController) UpdateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	id := c.Param("taskId") // パスパラメーター
	taskId, _ := strconv.Atoi(id)

	task := model.Task{}

	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	taskRes, err := tc.tu.UpdateTask(task, uint(userId.(float64)), uint(taskId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, taskRes)
}

func (tc *taskController) DeleteTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	id := c.Param("taskId")
	taskId, _ := strconv.Atoi(id)

	if err := tc.tu.DeleteTask(uint(userId.(float64)), uint(taskId)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)

}
