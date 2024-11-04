package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
	"strings"

	"github.com/go-vgo/robotgo"
)

func DoOrderToEarn() error {
	fmt.Println("DoOrderToEarn")

	allDone := true
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			return err
		}

		if ExistText("下单额外返元宝") {
			for _, v := range ocr.OCRResult {
				txt := v.([]interface{})[1].([]interface{})[0].(string)
				Polygon := v.([]interface{})[0]
				var txtRT, closeBtnLT, closeBtnRB robotgo.Point
				if txt == "下单额外返元宝" {
					txtRT.X = int(Polygon.([]interface{})[1].([]interface{})[0].(float64))
					txtRT.Y = int(Polygon.([]interface{})[1].([]interface{})[1].(float64))
					closeBtnLT.X = txtRT.X + (1204 - 1182)
					closeBtnRB.X = closeBtnLT.X + 16
					closeBtnLT.Y = txtRT.Y - (390 - 320)
					closeBtnRB.Y = closeBtnLT.Y + 16
					fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
					MoveClickTitle(closeBtnLT, closeBtnRB)
					robotgo.Sleep(2)
				}
			}
		}

		if OCRMoveClickTitle("立即签到", 0) {
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				return err
			}
			// 睡醒
			if OCRMoveClickTitle("我知道了", 0) {
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					return err
				}
			}
		}

		OCRMoveClickTitle("赚更多元宝", 0)

		for {
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				return err
			}

			var taskTitleLT, taskTitleRB robotgo.Point
			for _, v := range ocr.OCRResult {
				txt := v.([]interface{})[1].([]interface{})[0].(string)
				Polygon := v.([]interface{})[0]
				if strings.Contains(txt, "秒") || strings.Contains(txt, "分钟") {
					taskTitleLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
					taskTitleRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
					allDone = false
					break
				}
			}

			if allDone {
				break
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
					if todoBtnLT.Y > taskTitleLT.Y && taskTitleLT.Y < taskTitleRB.Y {
						MoveClickTitle(todoBtnLT, todoBtnRB)
						robotgo.Sleep(2)
						WatchAD("做任务赚元宝", "赚更多元宝")
						break
					}
				}
			}
		}

		newX := ocr.AppX + 18/2 + Utils.R.Intn(14/2)
		newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
		fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
		robotgo.MoveClick(newX, newY)
		robotgo.Sleep(2)

		if allDone {
			break
		}
	}

	newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)

	return nil
}