package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"

	"github.com/go-vgo/robotgo"
)

func DoDailySignIn() error {
	fmt.Println("DoDailySignIn")
loop:
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		// 20点抽奖
		if ExistText("立即抽奖") && OCRMoveClickTitle("立即抽奖", 0) {
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				panic(err)
			}

			if ExistText("做任务赚更多元宝") {
				OCRMoveClickTitle("做任务赚更多元宝", 0)
			}
		} else if ExistText("做任务赚更多元宝") { // 每日签到
			OCRMoveClickTitle("做任务赚更多元宝", 0)
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				panic(err)
			}
		} else if ExistText("继续做任务") { // 走到指定步数获奖
			OCRMoveClickTitle("继续做任务", 0)
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				panic(err)
			}
		}

		bNoTodo := true
		bDone := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0]
			if txt == "去完成" {
				bNoTodo = false
				break loop
			} else if txt == "已完成" {
				bDone = true
			}
		}

		// 没有一个 去完成 但有 已完成，那就是全部完成
		if bNoTodo && bDone {
			goto BACKTOINGOTCENTER
		}

		// 从下往上滑动
		newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
		newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
		robotgo.Move(newX, newY)
		fmt.Println("上滑 寻找 去完成")
		robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 90), 3, 50, Utils.R.Intn(10)-5)
		robotgo.Sleep(1)
	}

	fmt.Println("去完成 任务")
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if !ExistText("去完成") {
			break
		}

		OCRMoveClickTitle("去完成", 0)

		err = ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		WatchAD("签到赢元宝", "今日签到")
	}

BACKTOINGOTCENTER:
	newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)
	return nil
}
