package TaoliveOp

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

type TaskItem struct {
	// Title                string
	TitleLT, TitleRB     robotgo.Point
	TodoBtnLT, TodoBtnRB robotgo.Point
	Done                 bool
}

func GetTodoTask(taskList map[string]*TaskItem) *TaskItem {
	for title, taskItem := range taskList {
		fmt.Printf("%s: %v\n", title, taskItem.Done)
		if !taskItem.Done {
			return taskItem
		}
	}
	return nil
}
