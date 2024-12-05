package TaoliveOp

import (
	"fmt"
	"strings"

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
			if strings.Contains(taskItem.Title, "看黄金档直播并下单") && len(taskList) == 1 || !strings.Contains(taskItem.Title, "看黄金档直播并下单") {
				return taskItem
			}
		}
	}
	return nil
}
