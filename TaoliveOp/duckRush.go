package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"

	"github.com/go-vgo/robotgo"
)

func DoDuckRush() error {
	fmt.Println("DoDuckRush")

	// err := processWork()
	// if err != nil {
	// 	return err
	// }

	// x := ocr.AppX
	// y := ocr.AppY + ocr.AppHeight/2
	// w := ocr.AppWidth
	// h := ocr.AppHeight / 2
	// err := ocr.Ocr(&x, &y, &w, &h)
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return err
	}

	OCRMoveClickTitle(`^开始游戏$`, 0, true)

	for {
		x := ocr.AppX
		y := ocr.AppY + ocr.AppHeight/2
		w := ocr.AppWidth
		h := ocr.AppHeight / 2
		err := ocr.Ocr(&x, &y, &w, &h)
		if err != nil {
			return err
		}

		OCRMoveClickTitle(`^获得道具$`, 0, true)

		watchScroll10sAD()

		break
	}

	newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + 52/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)

	// newX = ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	// newY = ocr.AppY + 52/2 + Utils.R.Intn(26/2)
	// fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	// robotgo.MoveClick(newX, newY)
	// robotgo.Sleep(2)

	return nil
}
