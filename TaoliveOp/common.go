package TaoliveOp

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

type TaskItem struct {
	Title string
	// TitleLT, TitleRB     robotgo.Point
	TodoBtnLT, TodoBtnRB robotgo.Point
	// Done                 bool
}

func GetTodoTask(taskList []*TaskItem) *TaskItem {
	for _, taskItem := range taskList {
		fmt.Printf("(%3d, %3d)-(%3d, %3d): %v\n", taskItem.TodoBtnLT.X, taskItem.TodoBtnLT.Y, taskItem.TodoBtnRB.Y, taskItem.TodoBtnRB.Y, taskItem.Title)
		if taskItem.Title != "" {
			return taskItem
		}
	}
	return nil
}
