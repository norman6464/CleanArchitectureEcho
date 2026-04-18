package validator

import (
	"go-rest-api/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// バリデーター自体は他に依存をしないこと
// 依存をしているかしていないかはコンストラクタインジェクションを見ればわかる

type ITaskValidator interface {
	TaskValidate(task model.Task) error // バリデーションを行うためのメソッド
}

type taskValidator struct {
}

// 引数に何もないので依存先がない
func NewTaskValidator() ITaskValidator {
	return &taskValidator{}
}

func (tv *taskValidator) TaskValidate(task model.Task) error {
	return validation.ValidateStruct(&task,
		validation.Field(
			&task.Title,
			validation.Required.Error("title is required"),
			validation.RuneLength(1, 10).Error("limited max 10 char"),
		),
	)
}
