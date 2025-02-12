package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
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
	namedTask := make([]string, 0, 1)
	for _, taskItem := range taskList {
		if taskItem.Title != "" {
			namedTask = append(namedTask, taskItem.Title)
		}
	}

	for _, taskItem := range taskList {
		fmt.Printf("(%3d, %3d)-(%3d, %3d): %v\n", taskItem.TodoBtnLT.X, taskItem.TodoBtnLT.Y, taskItem.TodoBtnRB.Y, taskItem.TodoBtnRB.Y, taskItem.Title)
		if taskItem.Title != "" {
			if strings.Contains(taskItem.Title, "看黄金档直播并下单") && len(namedTask) == 1 || !strings.Contains(taskItem.Title, "看黄金档直播并下单") {
				return taskItem
			}
		}
	}
	return nil
}

func CloseAppStoreAD() {
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		if strings.Contains(txt, "给个好评") {
			Polygon := v.([]interface{})[0]
			var rightTop, closeBtnLT, closeBtnRB robotgo.Point
			rightTop.X = int(Polygon.([]interface{})[1].([]interface{})[0].(float64))
			rightTop.Y = int(Polygon.([]interface{})[1].([]interface{})[1].(float64))
			closeBtnLT.X = rightTop.X + (1250-1192)
			closeBtnLT.Y = rightTop.Y - (792-370)
			closeBtnRB.X = closeBtnLT.X + 22
			closeBtnRB.Y = closeBtnLT.Y + 22
			fmt.Printf("点击 关闭App Store评价(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
			MoveClickTitle(closeBtnLT, closeBtnRB)
			robotgo.Sleep(2)
			break
		}
	}
}