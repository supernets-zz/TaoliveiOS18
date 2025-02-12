package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"fmt"
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
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if !ExistText("去完成") && !ExistText("领元宝") {
			break
		}

		if !OCRMoveClickTitle(`^领元宝$`, 0, true) {
			OCRMoveClickTitle(`^去完成$`, 0, true)
		}

		err = ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		WatchAD("看短剧赚元宝", "看剧赚元宝")
	}

BACKTOSHORTVIDEOS:
	return fmt.Errorf("短剧任务已完成")
}
