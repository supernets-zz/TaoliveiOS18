package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

func DoShakeToEarn() error {
	fmt.Println("DoShakeToEarn")

	// err := processBrowseGetChances()
	// if err != nil {
	// 	return err
	// }
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return err
	}

	OCRMoveClickTitle(`^赚次数$`, 0, true)

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
			if txt == "去完成" {
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
		WatchAD("做任务赚次数", "赚次数")
	}

	newX := ocr.AppX + 18/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)

	newX = ocr.AppX + 18/2 + Utils.R.Intn(14/2)
	newY = ocr.AppY + 52/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)

	return nil
}

func processBrowseGetChances() error {
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return err
	}

	waitForEnter("赚次数", "")
	// for ExistText("立即领奖") || ExistText("视频福利") {
	// 	if OCRMoveClickTitle("立即领奖", 0) {
	// 		WatchAD("赚次数", "")
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
	// 			WatchAD("赚次数", "")
	// 			err := ocr.Ocr(nil, nil, nil, nil)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// }

slideLoop:
	for !ExistText("今日浏览任务已完成") {
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			re := regexp.MustCompile(`[（(]\d+/\d+[）)]`)
			matches := re.FindStringSubmatch(txt)
			if len(matches) > 0 {
				// 使用子表达式匹配
				rePart := regexp.MustCompile(`\d+`)
				parts := rePart.FindAllString(matches[0], -1)
				fmt.Println(parts) // 输出: [0 6]
				total, _ := strconv.ParseInt(parts[1], 10, 64)
				cur, _ := strconv.ParseInt(parts[0], 10, 64)
				if cur >= total {
					break slideLoop
				}

				// 先一滑到位，然后上下拉结合，上滑幅度比下拉大
				newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth*3/4)
				newY := ocr.AppY + ocr.AppHeight*3/4 + Utils.R.Intn(ocr.AppHeight/4)
				robotgo.Move(newX, newY)
				fmt.Println("上滑")
				robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 300), 6, 50, 0)
				downCnt := 1
				for i := 0; i < int(total-cur); i++ {
					for j := 0; j < 12; j++ {
						upOrDown := Utils.R.Intn(2)
						if j == 0 || !(upOrDown == 0 && downCnt > 0) {
							newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth*3/4)
							newY := ocr.AppY + ocr.AppHeight*3/4 + Utils.R.Intn(ocr.AppHeight/4)
							robotgo.Move(newX, newY)
							fmt.Println("上滑", downCnt)
							robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, 0)
							robotgo.Sleep(3)
							downCnt = downCnt + 1
						} else {
							newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth*3/4)
							newY := ocr.AppY + ocr.AppHeight/4 + Utils.R.Intn(ocr.AppHeight/4)
							robotgo.Move(newX, newY)
							fmt.Println("下拉", downCnt)
							robotgo.ScrollSmooth(-(Utils.R.Intn(30) - 100), 3, 50, 0)
							robotgo.Sleep(3)
							downCnt = downCnt - 1
						}
					}
				}

				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	// 上滑到顶
	for !ExistText("赚次数") {
		newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth*3/4)
		newY := ocr.AppY + ocr.AppHeight/4 + Utils.R.Intn(ocr.AppHeight/4)
		robotgo.Move(newX, newY)
		fmt.Println("下拉")
		robotgo.ScrollSmooth(-(Utils.R.Intn(30) - 400), 6, 50, 0)

		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
