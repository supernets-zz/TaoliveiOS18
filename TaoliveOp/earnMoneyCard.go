package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
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
		if OCRMoveClickTitle(`^领奖$`, 0, true) || OCRMoveClickTitle(`^看小视频30秒$`, 0, true) ||
			OCRMoveClickTitle(`^看精选推荐30秒$`, 0, true) || OCRMoveClickTitle(`^看上新好物30秒$`, 0, true) ||
			OCRMoveClickTitle(`^看省钱专区30秒$`, 0, true) {
			WatchAD("赚钱卡", "")
		} else if ContainText("可提现") {
			OCRMoveClickTitle(`^赚钱卡$`, 0, false)
		} else {
			break
		}
	}

	for {
		taskList := make([]*TaskItem, 0, 1)
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
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
				ti := &TaskItem{TodoBtnLT: todoBtnLT, TodoBtnRB: todoBtnRB}
				taskList = append(taskList, ti)
			}
		}

		var taskTitleLT, taskTitleRB robotgo.Point
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			Polygon := v.([]interface{})[0]
			if !strings.Contains(txt, "3元带走3件") && !strings.Contains(txt, "搜索并带走喜欢的宝贝") && !strings.Contains(txt, "Q") && !strings.Contains(txt, "3元") && txt != "已完成" && txt != "去完成" && txt != "领元宝" {
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
	newY := ocr.AppY + (414-306)/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)

	for {
		x := ocr.AppX
		y := ocr.AppY
		w := ocr.AppWidth
		h := ocr.AppHeight / 10
		err := ocr.Ocr(&x, &y, &w, &h)
		if err != nil {
			panic(err)
		}

		if ExistText("元宝中心") && ExistText("规则") {
			break
		}

		newX := ocr.AppX + 30/2 + Utils.R.Intn(14/2)
		newY := ocr.AppY + (414-306)/2 + Utils.R.Intn(26/2)
		fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
		robotgo.MoveClick(newX, newY)
		robotgo.Sleep(2)
	}

	return nil
}
