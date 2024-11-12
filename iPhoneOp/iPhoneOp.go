package IPhoneOp

import (
	ocr "TaoliveiOS18/OCR"
	"fmt"

	"github.com/andybrewer/mack"
	"github.com/go-vgo/robotgo"
)

func ActivateIPhoneMirroring() (int, error) {
	fpid, err := robotgo.FindIds("iPhone Mirroring")
	if err != nil {
		return 0, err
	}
	// fmt.Println("pids... ", fpid)

	robotgo.ActivePid(fpid[0])

	return fpid[0], nil
}

func GetIPhoneMirroringGeometry() error {
	fpid, err := ActivateIPhoneMirroring()
	if err != nil {
		return err
	}

	response, err := mack.Tell("System Events", fmt.Sprintf("set _P to a reference to (processes whose unix id is %d)", fpid), "set _W to a reference to windows of _P", "[_P's name, _W's size, _W's position]")
	if err != nil {
		return err
	}

	fmt.Println(response)
	var appName1, appName2 string
	_, err = fmt.Sscanf(response, "%s %s %d, %d, %d, %d", &appName1, &appName2, &ocr.AppWidth, &ocr.AppHeight, &ocr.AppX, &ocr.AppY)
	if err != nil {
		return err
	}
	fmt.Printf("原始(x, y, w, h) = (%d, %d, %d, %d)\n", ocr.AppX, ocr.AppY, ocr.AppWidth, ocr.AppHeight)
	ocr.AppX = ocr.AppX + 6
	ocr.AppY = ocr.AppY + 38
	ocr.AppWidth = ocr.AppWidth - 12
	ocr.AppHeight = ocr.AppHeight - 38 - 5
	fmt.Printf("修正后(x, y, w, h) = (%d, %d, %d, %d)\n", ocr.AppX, ocr.AppY, ocr.AppWidth, ocr.AppHeight)

	return nil
}

func GotoTaoLive() error {
	robotgo.KeyTap("1", "cmd")
	robotgo.KeyTap("1", "cmd")
	robotgo.Sleep(1)

	bNeedScroll := true
	for bNeedScroll {
		if err := ocr.Ocr(nil, nil, nil, nil); err != nil {
			panic(err)
		}

		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0]
			if txt == "点淘" {
				Polygon := v.([]interface{})[0]
				// fmt.Println(Polygon.([]interface{})[0].([]interface{})[0].(float64))
				var leftTop, rightBtm robotgo.Point
				leftTop.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
				leftTop.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
				rightBtm.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
				rightBtm.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
				// fmt.Println(x, y)
				// 点击 点淘 中心点
				fmt.Printf("点击 点淘(%3d, %3d)\n", ocr.AppX+int((leftTop.X+int((rightBtm.X-leftTop.X)/2))/2), ocr.AppY+int((leftTop.Y+int((rightBtm.Y-leftTop.Y)/2))/2)-25)
				robotgo.MoveClick(ocr.AppX+int((leftTop.X+int((rightBtm.X-leftTop.X)/2))/2), ocr.AppY+int((leftTop.Y+int((rightBtm.Y-leftTop.Y)/2))/2)-25)
				robotgo.Sleep(2)
				bNeedScroll = false
				break
			}
		}

		if bNeedScroll {
			// 从右往左小幅度滑动
			robotgo.Move(ocr.AppX+int(ocr.AppWidth*4/5), ocr.AppY+int(ocr.AppHeight/2))
			robotgo.ScrollSmooth(0, 3, 50, -300)
		}
	}

	return nil
}

func KillAllApp() error {

	return nil
}
