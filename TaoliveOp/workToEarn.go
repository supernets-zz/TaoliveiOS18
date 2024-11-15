package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-vgo/robotgo"
)

func DoWorkToEarn() error {
	fmt.Println("DoWorkToEarn")

	reTitle := regexp.MustCompile(`^(.*?)[\(（].*?$`)
	err := processWork()
	if err != nil {
		return err
	}

	OCRMoveClickTitle("得体力", 0)

	for {
		taskList := make(map[string]*TaskItem)
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			return err
		}

		var frameTitleRB robotgo.Point
		var taskTitleLT, taskTitleRB robotgo.Point
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			Polygon := v.([]interface{})[0]
			if txt == "做任务赚体力" {
				frameTitleRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			} else if strings.Contains(txt, "秒") || strings.Contains(txt, "分钟") {
				taskTitleLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
				taskTitleRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
				if frameTitleRB.Y > 0 && taskTitleRB.Y > frameTitleRB.Y {
					match := reTitle.FindStringSubmatch(txt)
					if len(match) > 1 {
						ti := &TaskItem{TitleLT: taskTitleLT, TitleRB: taskTitleRB, TodoBtnLT: taskTitleLT, TodoBtnRB: taskTitleLT, Done: false}
						taskList[match[1]] = ti
					}
				}
			}
		}

		var todoBtnLT, todoBtnRB robotgo.Point
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			Polygon := v.([]interface{})[0]
			if txt == "去完成" || txt == "已完成" {
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
		WatchAD("做任务赚体力", "得体力")
	}

	newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)

	newX = ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	newY = ocr.AppY + 52/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)

	return nil
}

func processWork() error {
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return err
	}

	waitForEnter("得体力", "")
	// for ExistText("立即领奖") || ExistText("视频福利") {
	// 	if OCRMoveClickTitle("立即领奖", 0) {
	// 		WatchAD("得体力", "")
	// 		err := ocr.Ocr(nil, nil, nil, nil)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else if OCRMoveClickTitle("视频福利", 0) {
	// 		err := ocr.Ocr(nil, nil, nil, nil)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		if OCRMoveClickTitle("立即领奖", 0) {
	// 			WatchAD("得体力", "")
	// 			err := ocr.Ocr(nil, nil, nil, nil)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// }

	if ExistText("领体力") && OCRMoveClickTitle("领体力", 50) {
		if OCRMoveClickTitle("放弃领取额外奖励", 0) {
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				return err
			}
		}
	}

	if ExistText("领取888元宝") && OCRMoveClickTitle("领取888元宝", 0) ||
		ExistText("领取588元宝") && OCRMoveClickTitle("领取588元宝", 0) ||
		ExistText("领取188元宝") && OCRMoveClickTitle("领取188元宝", 0) {
		if OCRMoveClickTitle("再得68元宝", 0) {
			WatchAD("得体力", "")
			robotgo.Sleep(3)
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				return err
			}
		}
	}

	if OCRMoveClickTitle("去打工赚钱", 0) {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			return err
		}

		if OCRMoveClickTitle("打工120分钟", 0) {
			OCRMoveClickTitle("开始打工", 0)
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
