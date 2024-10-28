package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"

	"github.com/go-vgo/robotgo"
)

func IsInTaoliveHome() (bool, error) {
	home, err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return false, err
	}

	var bLive, bStore, bShoppingCart bool
	for _, v := range home {
		txt := v.([]interface{})[1].([]interface{})[0]
		if txt == "直播" {
			bLive = true
			if bLive && bStore && bShoppingCart {
				return true, nil
			}
		} else if txt == "商城" {
			bStore = true
			if bLive && bStore && bShoppingCart {
				return true, nil
			}
		} else if txt == "购物车" {
			bShoppingCart = true
			if bLive && bStore && bShoppingCart {
				return true, nil
			}
		}
	}

	return bLive && bStore && bShoppingCart, nil
}

func IsInIngotCenter() (bool, error) {
	ingotCenter, err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return false, err
	}

	var bIngotCenter, bRule bool
	for _, v := range ingotCenter {
		txt := v.([]interface{})[1].([]interface{})[0]
		if txt == "元宝中心" {
			bIngotCenter = true
		} else if txt == "规则" {
			bRule = true
		}

		if bIngotCenter && bRule {
			return true, nil
		}
	}

	return bIngotCenter && bRule, nil
}

// iconHeight为文字上方图标的原分辨率高度
func MoveClickTitle(lt, rb robotgo.Point, iconHeight int) {
	// 截图是原分辨率，robotgo.MoveClick在Retina屏幕需要除以2
	robotgo.MoveClick(ocr.AppX+int((lt.X+Utils.R.Intn(rb.X-lt.X))/2), ocr.AppY+38+int((lt.Y-iconHeight/2+Utils.R.Intn(rb.Y-(lt.Y-iconHeight)))/2))
}

func GotoIngotCenter() error {
	ingotCenter, err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return err
	}

	for _, v := range ingotCenter {
		txt := v.([]interface{})[1].([]interface{})[0]
		if txt == "元宝中心" {
			Polygon := v.([]interface{})[0]
			// fmt.Println(Polygon.([]interface{})[0].([]interface{})[0].(float64))
			var leftTop, rightBtm robotgo.Point
			leftTop.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
			leftTop.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
			rightBtm.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
			rightBtm.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			// fmt.Println(x, y)
			fmt.Println("点击 元宝中心")
			MoveClickTitle(leftTop, rightBtm, 56)
			break
		}
	}

	return nil
}

func GotoDailySignIn() error {
	bNeedScroll := true
	var dailySignInLT, dailySignInRB robotgo.Point
	for bNeedScroll {
		dailySignIn, err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for _, v := range dailySignIn {
			txt := v.([]interface{})[1].([]interface{})[0]
			if txt == "今日签到" {
				Polygon := v.([]interface{})[0]
				// fmt.Println(Polygon.([]interface{})[0].([]interface{})[0].(float64))
				dailySignInLT.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
				dailySignInLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
				dailySignInRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
				dailySignInRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
				// fmt.Println(x, y)
				// 点击 今日签到
				fmt.Println("点击 今日签到")
				MoveClickTitle(dailySignInLT, dailySignInRB, 0)
				bNeedScroll = false
				break
			}
		}

		if bNeedScroll {
			// 从下往上滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + 38 + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 150), 3, 50, Utils.R.Intn(10)-5)
			robotgo.Sleep(1)
		}
	}

	return nil
}

// func ScrollInIngotCenter() error {
// 	newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
// 	newY := ocr.AppY + 38 + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
// 	robotgo.Move(newX, newY)
// 	robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 90), 3, 100, Utils.R.Intn(10)-5)

// 	return nil
// }
