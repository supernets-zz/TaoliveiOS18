package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
	"regexp"

	"github.com/go-vgo/robotgo"
)

func IsInTaoliveHome() (bool, error) {
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return false, err
	}

	var bLive, bStore, bShoppingCart, bIngotCenter, bShortVideos bool
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0]
		if txt == "直播" {
			bLive = true
			if bLive && bStore && bShoppingCart && bIngotCenter && bShortVideos {
				return true, nil
			}
		} else if txt == "商城" {
			bStore = true
			if bLive && bStore && bShoppingCart && bIngotCenter && bShortVideos {
				return true, nil
			}
		} else if txt == "购物车" {
			bShoppingCart = true
			if bLive && bStore && bShoppingCart && bIngotCenter && bShortVideos {
				return true, nil
			}
		} else if txt == "元宝中心" {
			bIngotCenter = true
			if bLive && bStore && bShoppingCart && bIngotCenter && bShortVideos {
				return true, nil
			}
		} else if txt == "短剧" {
			bShortVideos = true
			if bLive && bStore && bShoppingCart && bIngotCenter && bShortVideos {
				return true, nil
			}
		}
	}

	return bLive && bStore && bShoppingCart, nil
}

func IsInIngotCenter() (bool, error) {
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return false, err
	}

	var bIngotCenter, bRule bool
	for _, v := range ocr.OCRResult {
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

func IsInShortVideos() (bool, error) {
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return false, err
	}

	var bEarnIngot bool
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0]
		if txt == "看剧赚元宝" {
			bEarnIngot = true
		}

		if bEarnIngot {
			return true, nil
		}
	}

	return bEarnIngot, nil
}

func MoveClickTitle(leftTop, rightBtm robotgo.Point) {
	// 截图是原分辨率，robotgo.MoveClick在Retina屏幕需要除以2
	fmt.Printf("(%3d, %3d)-(%3d, %3d)\n", leftTop.X, leftTop.Y, rightBtm.X, rightBtm.Y)
	x := ocr.AppX + int((leftTop.X+Utils.R.Intn(rightBtm.X-leftTop.X))/2)
	y := ocr.AppY + int((leftTop.Y+Utils.R.Intn(rightBtm.Y-leftTop.Y))/2)
	robotgo.MoveClick(x, y)
	fmt.Printf("点击 (%3d, %3d)\n", x, y)
}

// iconHeight为文字上方图标的原分辨率高度
func OCRMoveClickTitle(pattern string, iconHeight int, untilChange bool) bool {
	fmt.Printf("-> OCRMoveClickTitle(%s)\n", pattern)
	defer fmt.Printf("<- OCRMoveClickTitle(%s)\n", pattern)
	bClick := false
	// 截图是原分辨率，robotgo.MoveClick在Retina屏幕需要除以2
	var leftTop, rightBtm robotgo.Point
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		matched, err := regexp.Match(pattern, []byte(txt))
		if err != nil {
			fmt.Println(err)
			break
		}
		if matched {
			Polygon := v.([]interface{})[0]
			// fmt.Println(Polygon.([]interface{})[0].([]interface{})[0].(float64))
			leftTop.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
			leftTop.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64)) - int(iconHeight/2)
			rightBtm.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
			rightBtm.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			// x := ocr.AppX + int((leftTop.X+Utils.R.Intn(rightBtm.X-leftTop.X))/2)
			// y := ocr.AppY + int((leftTop.Y-iconHeight/2+Utils.R.Intn(rightBtm.Y-(leftTop.Y-iconHeight)))/2)
			x := ocr.AppX + int((leftTop.X+Utils.R.Intn(rightBtm.X-leftTop.X))/2)
			y := ocr.AppY + int((leftTop.Y+Utils.R.Intn(rightBtm.Y-leftTop.Y))/2)
			// 点击 去完成
			fmt.Printf("点击 %s(%3d, %3d)\n", pattern, x, y)
			robotgo.MoveClick(x, y)
			robotgo.Sleep(2)
			robotgo.Move(ocr.AppX+ocr.AppWidth+10, ocr.AppY+ocr.AppHeight+10)
			bClick = true
			break
		}
	}

	if bClick && untilChange {
		x := ocr.AppX + leftTop.X/2 - 10
		y := ocr.AppY + leftTop.Y/2 - 10
		w := (rightBtm.X-leftTop.X)/2 + 20
		h := (rightBtm.Y-leftTop.Y)/2 + 20
		err := ocr.Ocr(&x, &y, &w, &h)
		if err != nil {
			panic(err)
		}

		for {
			bClickSucc := true
			for _, v := range ocr.OCRResult {
				txt := v.([]interface{})[1].([]interface{})[0].(string)
				matched, err := regexp.Match(pattern, []byte(txt))
				if err != nil {
					fmt.Println(err)
					break
				}
				if matched {
					bClickSucc = false
					break
				}
			}

			if bClickSucc {
				robotgo.Sleep(2)
				break
			}

			x := ocr.AppX + int((leftTop.X+Utils.R.Intn(rightBtm.X-leftTop.X))/2)
			y := ocr.AppY + int((leftTop.Y+Utils.R.Intn(rightBtm.Y-leftTop.Y))/2)
			// 点击 去完成
			fmt.Printf("点击 %s(%3d, %3d)\n", pattern, x, y)
			robotgo.MoveClick(x, y)
			robotgo.Sleep(5)

			x = ocr.AppX + leftTop.X/2 - 10
			y = ocr.AppY + leftTop.Y/2 - 10
			w := (rightBtm.X-leftTop.X)/2 + 20
			h := (rightBtm.Y-leftTop.Y)/2 + 20
			err := ocr.Ocr(&x, &y, &w, &h)
			if err != nil {
				panic(err)
			}
		}
	}
	return bClick
}

func GotoIngotCenter() error {
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return err
	}

	for {
		bInIngotCenter, err := IsInIngotCenter()
		if err != nil {
			panic(err)
		}

		if bInIngotCenter {
			break
		}

		OCRMoveClickTitle(`^元宝中心$`, 0, true)
	}

	return nil
}

func GotoShortVideos() error {
	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		return err
	}

	for {
		bInIngotCenter, err := IsInShortVideos()
		if err != nil {
			panic(err)
		}

		if bInIngotCenter {
			break
		}

		OCRMoveClickTitle(`^看剧赚元宝$`, 0, true)
	}

	return nil
}

func GotoDailySignIn() error {
	bNeedScroll := true
	for bNeedScroll {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for ExistText("立即领奖") {
			if OCRMoveClickTitle(`^立即领奖$`, 0, true) {
				WatchAD("元宝中心", "")
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					panic(err)
				}
			}
		}

		if OCRMoveClickTitle(`^今日签到$`, 0, true) {
			bNeedScroll = false
			break
		}

		if bNeedScroll {
			// 从下往上滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 150), 3, 50, Utils.R.Intn(10)-5)
			robotgo.Sleep(1)
		}
	}

	return nil
}

func GotoEarnMoneyCard() error {
	bNeedScroll := true
	for bNeedScroll {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for ExistText("立即领奖") {
			if OCRMoveClickTitle(`^立即领奖$`, 0, true) {
				WatchAD("元宝中心", "")
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					panic(err)
				}
			}
		}

		if OCRMoveClickTitle(`^赚钱卡$`, 0, true) {
			bNeedScroll = false
			break
		}

		if bNeedScroll {
			// 从下往上滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 150), 3, 50, Utils.R.Intn(10)-5)
			robotgo.Sleep(1)
		}
	}

	return nil
}

func GotoWalkToEarn() error {
	bNeedScroll := true
	for bNeedScroll {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for ExistText("立即领奖") {
			if OCRMoveClickTitle(`^立即领奖$`, 0, true) {
				WatchAD("元宝中心", "")
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					panic(err)
				}
			}
		}

		if OCRMoveClickTitle(`^走路赚元宝$`, 0, true) {
			bNeedScroll = false
			break
		}

		if bNeedScroll {
			// 从下往上滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 150), 3, 50, Utils.R.Intn(10)-5)
			robotgo.Sleep(1)
		}
	}

	return nil
}

func GotoWorkToEarn() error {
	bNeedScroll := true
	for bNeedScroll {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for ExistText("立即领奖") {
			if OCRMoveClickTitle(`^立即领奖$`, 0, true) {
				WatchAD("元宝中心", "")
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					panic(err)
				}
			}
		}

		if OCRMoveClickTitle(`^打工赚元宝$`, 0, true) {
			bNeedScroll = false
			break
		}

		if bNeedScroll {
			// 从下往上滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 150), 3, 50, Utils.R.Intn(10)-5)
			robotgo.Sleep(1)
		}
	}

	return nil
}

func GotoShakeToEarn() error {
	bNeedScroll := true
	for bNeedScroll {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for ExistText("立即领奖") {
			if OCRMoveClickTitle("立即领奖", 0, true) {
				WatchAD("元宝中心", "")
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					panic(err)
				}
			}
		}

		if OCRMoveClickTitle("摇一摇赚元宝", 0, true) {
			bNeedScroll = false
			break
		}

		if bNeedScroll {
			// 从下往上滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 150), 3, 50, Utils.R.Intn(10)-5)
			robotgo.Sleep(1)
		}
	}

	return nil
}

func GotoSleepToEarn() error {
	bNeedScroll := true
	for bNeedScroll {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for ExistText("立即领奖") {
			if OCRMoveClickTitle("立即领奖", 0, true) {
				WatchAD("元宝中心", "")
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					panic(err)
				}
			}
		}

		if OCRMoveClickTitle("睡觉赚元宝", 0, true) {
			bNeedScroll = false
			break
		}

		if bNeedScroll {
			// 从下往上滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 150), 3, 50, Utils.R.Intn(10)-5)
			robotgo.Sleep(1)
		}
	}

	return nil
}

func GotoOrderToEarn() error {
	bNeedScroll := true
	for bNeedScroll {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for ExistText("立即领奖") {
			if OCRMoveClickTitle("立即领奖", 0, true) {
				WatchAD("元宝中心", "")
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					panic(err)
				}
			}
		}

		if OCRMoveClickTitle("下单返元宝", 0, true) {
			bNeedScroll = false
			break
		}

		if bNeedScroll {
			// 从上往下滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/4 + Utils.R.Intn(ocr.AppHeight/4)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(10) - 100), 3, 50, 0)
			robotgo.Sleep(1)
		}
	}

	return nil
}

func GotoDuckRush() error {
	bNeedScroll := true
	for bNeedScroll {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if OCRMoveClickTitle(`^鸭鸭快跑$`, 0, true) {
			bNeedScroll = false
			break
		}

		if bNeedScroll {
			// 从下往上滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
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
