package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
)

// controller型が参照（依存）するので定義をする
type ITaskUsecase interface {
	GetAllTasks(userId uint) ([]model.TaskResponse, error)
	GetTaskById(userId uint, taskId uint) (model.TaskResponse, error)
	CreateTask(task model.Task) (model.TaskResponse, error)
	UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error)
	DeleteTask(userId uint, taskId uint) error
}

type taskUsecase struct {
	tr repository.ITaskRepository
}

// コンストラクタ
func NewTaskUsecase(tr repository.ITaskRepository) ITaskUsecase {
	return &taskUsecase{tr} // インスタンス化
}
