package repository

import (
	"fmt"
	"go-rest-api/model"

	"gorm.io/gorm"
)

type ITaskRepository interface {
	GetAllTasks(tasks *[]model.Task, userId uint) error
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task, userId uint, taskId uint) error
	DeleteTask(userId uint, taskId uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db}
}

// repositoryの具体の実装をする
// 戻り値はないがこのリポジトリ（interface）を依存関係とするusecaseレイヤー側にアドレスで格納をさせる
func (tr *taskRepository) GetAllTasks(tasks *[]model.Task, userId uint) error {
	if err := tr.db.Joins("User").Where("user_id=?", userId).Order("created_at").Find(tasks).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	// Firstでは第一引数ではtaskフィールドの中のtaskIdが一致するものをselectする
	if err := tr.db.Joins("User").Where("user_id=?", userId).First(task, taskId).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) CreateTask(task *model.Task) error {
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) UpdateTask(task *model.Task, userId uint, taskId uint) error {

}

func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
	result := tr.db.Where("id=? AND user_id=?", taskId, userId).Delete(&model.Task{})

	if result.Error != nil {
		return result.Error
	}

	// gorm.DB型のDelete()、update()の戻り値はRowAffectedで何行deleteやupdateされたのかを保存する
	if result.RowsAffected < 1 {
		return fmt.Errorf("Object does not exist")
	}
	return nil
}
