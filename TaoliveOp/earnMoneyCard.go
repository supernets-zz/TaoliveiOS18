package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"

	"github.com/go-vgo/robotgo"
)

func DoEarnMoneyCard() error {
	fmt.Println("DoEarnMoneyCard")
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if OCRMoveClickTitle("领奖", 0) || OCRMoveClickTitle("看小视频30秒", 0) ||
			OCRMoveClickTitle("看精选推荐30秒", 0) || OCRMoveClickTitle("看上新好物30秒", 0) ||
			OCRMoveClickTitle("看省钱专区30秒", 0) {
			WatchAD("赚钱卡")
		} else if containText("可提现") {
			OCRMoveClickTitle("赚钱卡", 0)
		} else {
			break
		}
	}

loop2:
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		bNoTodo := true
		bDone := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0]
			if txt == "去完成" || txt == "领元宝" {
				bNoTodo = false
				break loop2
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
		robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 90), 3, 50, Utils.R.Intn(10)-5)
		robotgo.Sleep(1)
	}

	fmt.Println("去完成 任务")
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if !OCRMoveClickTitle("领元宝", 0) && !OCRMoveClickTitle("去完成", 0) {
			break
		}

		err = ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		WatchAD("赚钱卡")
	}

BACKTOINGOTCENTER:
	newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)
	return nil
}
