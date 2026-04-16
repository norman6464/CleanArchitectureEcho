package repository

import (
	"fmt"
	"go-rest-api/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	// Clausesメソッドはgormにupdate実行後DBからの更新後の行をselectして構造体に書き直す用の指示する句
	// Clausesを実行しなかった場合、構造体taskに変更後の反映がされない
	result := tr.db.Model(task).Clauses(clause.Returning{}).Where("id=? AND user_id=?", taskId, userId).Update("title", task.Title)
	if result.Error != nil {
		return result.Error
	}
	// Update文は戻り値はRowAffectedで何行updateされたのかをRowsAffectedで確認ができる
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exsit")
	}

	return nil

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
