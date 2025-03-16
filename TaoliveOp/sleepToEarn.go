package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"

	"github.com/go-vgo/robotgo"
)

func DoSleepToEarn() error {
	fmt.Println("DoSleepToEarn")
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return err
	}

	waitForEnter("定提醒", "")
	// for ExistText("立即领奖") || ExistText("视频福利") || ExistText("看5s得") {
	// 	if OCRMoveClickTitle("立即领奖", 0) {
	// 		WatchAD("定提醒", "")
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
	// 			WatchAD("定提醒", "")
	// 			err := ocr.Ocr(nil, nil, nil, nil)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 	} else if OCRMoveClickTitle("看5s得", 0) {
	// 		WatchAD("定提醒", "")
	// 		err := ocr.Ocr(nil, nil, nil, nil)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	if OCRMoveClickTitle(`^可得666元宝$`, 0, true) || OCRMoveClickTitle(`^可得133元宝$`, 0, true) || OCRMoveClickTitle(`^可得333元宝$`, 0, true) {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			return err
		}
		// 睡醒
		if OCRMoveClickTitle(`^浏览30秒至少得68元宝$`, 0, true) {
			WatchAD("定提醒", "")
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				return err
			}
		}
	}

	if !ExistText("+38") && !ExistText("+48") && !ExistText("+58") && !ExistText("+68") && !ExistText("+100") {
		newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
		newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
		fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
		robotgo.MoveClick(newX, newY)
		robotgo.Sleep(2)
		return nil
	}

	for {
		if OCRMoveClickTitle(`^.?38$`, 0, true) || OCRMoveClickTitle(`^.?48$`, 0, true) || OCRMoveClickTitle(`^.?58$`, 0, true) || OCRMoveClickTitle(`^.?68$`, 0, true) || OCRMoveClickTitle(`^.?100$`, 0, true) {
			WatchAD("定提醒", "")
		} else {
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
