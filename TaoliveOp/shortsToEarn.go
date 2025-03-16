package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"fmt"
	"strings"

	"github.com/go-vgo/robotgo"
)

func DoShortsToEarn() error {
	fmt.Println("DoShortsToEarn")
loop:
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		bNoTodo := true
		bDone := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			fmt.Println(txt)
			if txt == "去完成" || txt == "领元宝" {
				bNoTodo = false
				break loop
			} else if txt == "已完成" {
				bDone = true
			}
		}

		// 没有一个 去完成 但有 已完成，那就是全部完成
		if bNoTodo && bDone {
			goto BACKTOSHORTVIDEOS
		}
	}

	fmt.Println("去完成 任务")
	for {
		taskList := make([]*TaskItem, 0, 1)
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			return err
		}

		var todoBtnLT, todoBtnRB robotgo.Point
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			Polygon := v.([]interface{})[0]
			if txt == "去完成" || txt == "领元宝" {
				todoBtnLT.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
				todoBtnLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
				todoBtnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
				todoBtnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
				ti := &TaskItem{TodoBtnLT: todoBtnLT, TodoBtnRB: todoBtnRB}
				taskList = append(taskList, ti)
			}
		}

		var taskTitleLT, taskTitleRB robotgo.Point
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			Polygon := v.([]interface{})[0]
			if strings.Contains(txt, "秒") || strings.Contains(txt, "分钟") {
				taskTitleLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
				taskTitleRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
				for _, taskItem := range taskList {
					if taskItem.TodoBtnLT.Y > taskTitleLT.Y-5 && taskItem.TodoBtnLT.Y < taskTitleRB.Y+5 {
						fmt.Println(txt)
						taskItem.Title = txt
						break
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
		WatchAD("看短剧赚元宝", "看剧赚元宝")
	}

BACKTOSHORTVIDEOS:
	return fmt.Errorf("短剧任务已完成")
}
