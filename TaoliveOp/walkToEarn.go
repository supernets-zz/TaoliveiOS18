package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

var Steps = []string{"100步", "500步", "2500步", "5000步", "10000步", "12500步", "15000步", "17500步", "20000步"}
var curStepsTipsRB, canWalkStepsTipsRB robotgo.Point
var curSteps, canWalkSteps int64
var validCheckPoints []bool
var checkPoints []PtPair

// 100、500、2500、5000、10000、12500、15000、17500、20000
func DoWalkToEarn() error {
	fmt.Println("DoWalkToEarn")

	validCheckPoints = make([]bool, len(Steps))
	checkPoints = make([]PtPair, len(Steps))

	// 当前步数大于气泡的步数，点气泡
	err := processBubbles()
	if err != nil {
		return err
	}

	bTodayDone := ExistText("今日步数已完成") || ExistText("今日已走完")
	fmt.Printf("今日已走：%d, 今日步数已完成: %v\n", curSteps, bTodayDone)
	if curSteps >= 20000 || bTodayDone {
		newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
		newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
		fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
		robotgo.MoveClick(newX, newY)
		robotgo.Sleep(2)
		return nil
	}

	OCRMoveClickTitle(`^赚步数$`, 0, true)

	for i := 0; i < int(math.Min(10, math.Ceil(float64(20000-curSteps)/150))); i++ {
		taskList := make([]*TaskItem, 0, 1)
		err = ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			return err
		}

		var todoBtnLT, todoBtnRB robotgo.Point
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			Polygon := v.([]interface{})[0]
			if txt == "去完成" || txt == "去浏览" {
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
			if (strings.Contains(txt, "秒") || strings.Contains(txt, "分钟")) && !strings.Contains(txt, "元宝商城") {
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
		WatchAD("做任务赚步数", "赚步数")
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

func processBubbles() error {
	robotgo.Sleep(3)

	x := ocr.AppX
	y := ocr.AppY
	w := ocr.AppWidth
	h := int(ocr.AppHeight * 4 / 5)
	err := ocr.Ocr(&x, &y, &w, &h)
	if err != nil {
		return err
	}

	for ExistText("立即领奖") || ExistText("视频福利") {
		if OCRMoveClickTitle("立即领奖", 0, true) {
			WatchAD("赚步数", "")
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				return err
			}
		} else if OCRMoveClickTitle("视频福利", 0, true) {
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				return err
			}

			if OCRMoveClickTitle("立即领奖", 0, true) {
				WatchAD("赚步数", "")
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	if ExistText("领取") && OCRMoveClickTitle("领取", 50, true) {
		newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
		newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
		fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
		robotgo.MoveClick(newX, newY)
		robotgo.Sleep(2)
	}

	if ExistText("今日步数已完成") {
		return nil
	}

	for i := 0; i < len(validCheckPoints); i++ {
		validCheckPoints[i] = false
	}

	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		Polygon := v.([]interface{})[0]
		if txt == "当前步数" || txt == "今日已走" {
			curStepsTipsRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
		} else if txt == "可用步数" {
			canWalkStepsTipsRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
		} else {
			idx := Utils.IndexOf(Steps, txt)
			if idx != -1 {
				validCheckPoints[idx] = true
				checkPoints[idx].LeftTop.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
				checkPoints[idx].LeftTop.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
				checkPoints[idx].RightBtm.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
				checkPoints[idx].RightBtm.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			}
		}
	}

	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		Polygon := v.([]interface{})[0]
		var stepsLT robotgo.Point
		stepsLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
		if Utils.IsNumeric(txt) && stepsLT.Y-curStepsTipsRB.Y < 10 {
			curSteps, err = strconv.ParseInt(txt, 10, 64)
			if err != nil {
				return err
			}
		} else if Utils.IsNumeric(txt) && stepsLT.Y-canWalkStepsTipsRB.Y < 10 {
			canWalkSteps, err = strconv.ParseInt(txt, 10, 64)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("当前步数: %d, 可用步数: %d\n", curSteps, canWalkSteps)

	firstUnclickCheckPointIdx := -1
	for i := 0; i < len(Steps); i++ {
		if validCheckPoints[i] {
			if firstUnclickCheckPointIdx == -1 {
				firstUnclickCheckPointIdx = i
			}
			fmt.Printf("%s, (%3d, %3d)-(%3d, %3d)\n", Steps[i],
				checkPoints[i].LeftTop.X, checkPoints[i].LeftTop.Y,
				checkPoints[i].RightBtm.X, checkPoints[i].RightBtm.Y)
		}
	}

	if canWalkSteps > 0 {
		OCRMoveClickTitle("出发", 0, false)
		curSteps = curSteps + canWalkSteps
	}

	if firstUnclickCheckPointIdx != -1 {
		for i := firstUnclickCheckPointIdx; i < len(Steps); i++ {
			unclickStep, err := Utils.ExtractNumber(Steps[i])
			if err != nil {
				return err
			}
			if curSteps >= int64(unclickStep) {
				OCRMoveClickTitle(Steps[i], 34, true)
				robotgo.Sleep(2)
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					panic(err)
				}
				if OCRMoveClickTitle("浏览30秒再得68元宝", 0, true) || OCRMoveClickTitle("再得68元宝", 0, true) {
					WatchAD("赚步数", "")
				}
			} else {
				break
			}
		}
	}

	return nil
}
