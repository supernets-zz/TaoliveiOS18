package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-vgo/robotgo"
)

func DoEarnMoneyCard() error {
	fmt.Println("DoEarnMoneyCard")
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		waitForEnter("赚钱卡", "")
		if OCRMoveClickTitle("领奖", 0) || OCRMoveClickTitle("看小视频30秒", 0) ||
			OCRMoveClickTitle("看精选推荐30秒", 0) || OCRMoveClickTitle("看上新好物30秒", 0) ||
			OCRMoveClickTitle("看省钱专区30秒", 0) {
			WatchAD("赚钱卡", "")
		} else if containText("可提现") {
			OCRMoveClickTitle("赚钱卡", 0)
		} else {
			break
		}
	}

	reTitle := regexp.MustCompile(`^(.*?)[\(（].*?$`)
	for {
		taskList := make(map[string]*TaskItem)
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		var taskTitleLT, taskTitleRB robotgo.Point
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			Polygon := v.([]interface{})[0]
			if !strings.Contains(txt, "3元带走3件") {
				taskTitleLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
				taskTitleRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
				match := reTitle.FindStringSubmatch(txt)
				if len(match) > 1 {
					ti := &TaskItem{TitleLT: taskTitleLT, TitleRB: taskTitleRB, TodoBtnLT: taskTitleLT, TodoBtnRB: taskTitleLT, Done: false}
					taskList[match[1]] = ti
				}
			}
		}

		var todoBtnLT, todoBtnRB robotgo.Point
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0]
			Polygon := v.([]interface{})[0]
			if txt == "去完成" || txt == "领元宝" {
				todoBtnLT.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
				todoBtnLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
				todoBtnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
				todoBtnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
				for title, taskItem := range taskList {
					if todoBtnLT.Y > taskItem.TitleLT.Y-5 && todoBtnLT.Y < taskItem.TitleRB.Y+5 {
						fmt.Printf("%s: %s\n", title, txt)
						taskItem.TodoBtnLT.X = todoBtnLT.X
						taskItem.TodoBtnLT.Y = todoBtnLT.Y
						taskItem.TodoBtnRB.X = todoBtnRB.X
						taskItem.TodoBtnRB.Y = todoBtnRB.Y
						if txt == "去完成" {
							taskItem.Done = false
						} else if txt == "已完成" {
							taskItem.Done = true
						}
					}
				}
			}
		}

		taskItem := GetTodoTask(taskList)
		if taskItem == nil {
			break
		}

		MoveClickTitle(taskItem.TodoBtnLT, taskItem.TodoBtnRB)
		robotgo.Sleep(2)
		WatchAD("赚钱卡", "")

		// // 从下往上滑动
		// newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
		// newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
		// robotgo.Move(newX, newY)
		// robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 90), 3, 50, Utils.R.Intn(10)-5)
		// robotgo.Sleep(1)
	}

	// BACKTOINGOTCENTER:
	newX := ocr.AppX + 30/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + (246-138)/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)
	return nil
}
