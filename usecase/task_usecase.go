package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
)

// 基本的ユースケース層は今までrepositoryに依存していたがこれをインターフェースに依存をすることで依存性逆転の原則が適用できる

// 内側になるがインターフェースをさらに設けることでusecaseに依存をするcontroller側の変更容易性が高まる
type ITaskUsecase interface {
	GetAllTasks(userId uint) ([]model.TaskResponse, error)
	GetTaskById(userId uint, taskId uint) (model.TaskResponse, error)
	CreateTask(task model.Task) (model.TaskResponse, error)
	UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error)
	DeleteTask(userId uint, taskId uint) error
}

type taskUsecase struct {
	tr repository.ITaskRepository // 依存関係はインターフェースにしている
}

func NewTaskUsecase(tr repository.ITaskRepository) ITaskUsecase {
	//　インターフェースを戻り値の方にしているが、実際には構造体の値を返している
	// これができるのは構造体に紐づいている関数
	return &taskUsecase{tr}
}

func (tu *taskUsecase) GetAllTasks(userId uint) ([]model.TaskResponse, error) {
	tasks := []model.Task{} // Task構造体のリスト化
	if err := tu.tr.GetAllTasks(&tasks, userId); err != nil {
		return nil, err
	}

	resTasks := []model.TaskResponse{}
	for _, v := range tasks { // _はKey、vはValue
		t := model.TaskResponse{
			ID:        v.ID,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		resTasks = append(resTasks, t)
	}
	return resTasks, nil
}

func (tu *taskUsecase) GetTaskById(userId uint, taskId uint) (model.TaskResponse, error) {
	task := model.Task{}
	if err := tu.tr.GetTaskById(&task, userId, taskId); err != nil {
		return model.TaskResponse{}, err // からのインスタンスとerrオブジェクトを返す
	}
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}

func (tu *taskUsecase) CreateTask(task model.Task) (model.TaskResponse, error) {
	if err := tu.tr.CreateTask(&task); err != nil {
		return model.TaskResponse{}, err // 空のインスタンス化したものを返却をする
	}
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}

func (tu *taskUsecase) UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error) {
	if err := tu.tr.UpdateTask(&task, userId, taskId); err != nil {
		return model.TaskResponse{}, nil // エラーの時はまずは空のインスタンスを生成する
	}

	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil

}

func (tu *taskUsecase) DeleteTask(userId uint, taskId uint) error {
	if err := tu.tr.DeleteTask(userId, taskId); err != nil {
		return err
	}
	return nil
}
