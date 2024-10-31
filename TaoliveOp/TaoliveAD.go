package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

const (
	// 直播广告
	ADType_Live = iota
	// 每6秒向上滑动广告
	ADType_Scroll6s
	// 每10秒向上滑动广告
	ADType_Scroll10s
	// 点击广告后可跳过
	ADType_ClickADToSkip
	ADType_Unknown
)

type PtPair struct {
	LeftTop  robotgo.Point
	RightBtm robotgo.Point
}

var chADStart, chADEnd chan struct{}

var predictADEndTick int64

func ADType(result []interface{}) int {
	for _, v := range result {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		if txt == "更多直播" {
			return ADType_Live
		} else if strings.Contains(txt, "滑动浏览") {
			return ADType_Scroll6s
		} else if strings.Contains(txt, "点击广告可领取奖励|跳过") {
			return ADType_ClickADToSkip
		}
	}

	return ADType_Unknown
}

func ExistText(text string) bool {
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		if txt == text {
			return true
		}
	}

	return false
}

func containText(text string) bool {
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		if strings.Contains(txt, text) {
			return true
		}
	}

	return false
}

func waitForLeave(scene string) {
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if !ExistText(scene) {
			break
		}
	}
}

func waitForEnter(scene string) {
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if ExistText("继续做任务") {
			OCRMoveClickTitle("继续做任务", 0)
		}

		if ExistText(scene) {
			break
		}
	}
}

func monitor(title string, fn func() bool) {
	<-chADStart
	for {
		if predictADEndTick != 0 {
			fmt.Printf("%s, 预计广告结束时间: %s\n", title, time.Unix(predictADEndTick, 0).Local().Format("2006-01-02 15:04:05"))
		}
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if fn() {
			goto ADEnd
		}
		// for _, v := range ocr.OCRResult {
		// 	txt := v.([]interface{})[1].([]interface{})[0].(string)
		// 	if strings.Contains(txt, "滑动浏览") {
		// 		fmt.Println(txt)
		// 		goto ADEnd
		// 	}
		// }
	}
ADEnd:
	fmt.Println("广告结束")
	chADEnd <- struct{}{}
}

func WatchAD(lastScene string) {
	chADStart = make(chan struct{})
	chADEnd = make(chan struct{})
	defer close(chADStart)
	defer close(chADEnd)

	if lastScene != "" {
		fmt.Println("等待离开", lastScene)
		waitForLeave(lastScene)
	}

	if (containText("秒后完成") || containText("秒后领奖") || containText("秒后发奖")) && containText("说点") {
		watchLiveAD()
	} else if containText("秒后完成") && !containText("说点") {
		watchVideoAD()
	} else if containText("滑动浏览") {
		watchScroll6sAD()
	} else if containText("秒后发放奖励") || containText("计时已暂停，上滑继续") {
		watchScroll10sAD()
	} else if containText("点击广告可领取奖励|跳过") {
		watchClickToSkipAD()
	} else if containText("搜索领元宝") || containText("搜索领体力") {
		watchSearchScrollAD()
	} else if containText("立即获取") && containText("跳过") {
		watchInteractiveAD()
	} else if containText("s|跳过") {

	} else {
		fmt.Println("Unknown")
		watchLiveAD()
	}

	if lastScene != "" {
		fmt.Println("等待回到", lastScene)
		waitForEnter(lastScene)
	}
}

func watchLiveAD() {
	fmt.Println("watchLiveAD")
	predictADEndTick = 0
	go monitor("watchLiveAD", func() bool {
		bCountdownComplete := true
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if strings.Contains(txt, "秒后完成") {
				if predictADEndTick == 0 {
					leftSec, _ := Utils.ExtractNumber(txt)
					predictADEndTick = time.Now().Unix() + int64(leftSec+3)
				}
				bCountdownComplete = false
				break
			}
		}
		return bCountdownComplete && time.Now().Unix() > predictADEndTick
	})
	chADStart <- struct{}{}

	startTick := time.Now()
loop:
	for {
		if time.Now().Unix()-startTick.Unix() > 240 {
			fmt.Println("timeout")
			break loop
		}
		select {
		case <-time.After((10 + time.Duration(Utils.R.Intn(2))) * time.Second):
			if ExistText("领取") {
				OCRMoveClickTitle("领取", 0)
			}
			if ExistText("点击x2倍") {
				OCRMoveClickTitle("点击x2倍", 0)
			}
			if ExistText("点击x4倍") {
				OCRMoveClickTitle("点击x4倍", 0)
			}
			for _, v := range ocr.OCRResult {
				txt := v.([]interface{})[1].([]interface{})[0].(string)
				Polygon := v.([]interface{})[0]
				var btnLB, btnRB, closeBtnLT, closeBtnRB robotgo.Point
				if txt == "立即领取" {
					btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
					btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
					btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
					btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
					centerX := btnLB.X + int((btnRB.X-btnLB.X)/2)
					closeBtnLT.X = centerX - 22/2
					closeBtnRB.X = centerX + 22/2
					closeBtnLT.Y = btnLB.Y + (888 - 786)
					closeBtnRB.Y = closeBtnLT.Y + 22
					fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
					MoveClickTitle(closeBtnLT, closeBtnRB)
					robotgo.Sleep(2)
				} else if txt == "再来一次" {
					btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
					btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
					btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
					btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
					centerX := btnLB.X + int((btnRB.X-btnLB.X)/2)
					closeBtnLT.X = centerX - 22/2
					closeBtnRB.X = centerX + 22/2
					closeBtnLT.Y = btnLB.Y + (952 - 774)
					closeBtnRB.Y = closeBtnLT.Y + 22
					fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
					MoveClickTitle(closeBtnLT, closeBtnRB)
					robotgo.Sleep(2)
				} else if txt == "关注并领取" {
					btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
					btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
					btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
					btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
					centerX := btnLB.X + int((btnRB.X-btnLB.X)/2)
					closeBtnLT.X = centerX - 22/2
					closeBtnRB.X = centerX + 22/2
					closeBtnLT.Y = btnLB.Y + (952 - 774)
					closeBtnRB.Y = closeBtnLT.Y + 22
					fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
					MoveClickTitle(closeBtnLT, closeBtnRB)
					robotgo.Sleep(2)
				} else if txt == "开心收下" {
					btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
					btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
					btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
					btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
					centerX := btnLB.X + int((btnRB.X-btnLB.X)/2)
					closeBtnLT.X = centerX - 22/2
					closeBtnRB.X = centerX + 22/2
					closeBtnLT.Y = btnLB.Y + (952 - 774)
					closeBtnRB.Y = closeBtnLT.Y + 22
					fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
					MoveClickTitle(closeBtnLT, closeBtnRB)
					robotgo.Sleep(2)
				}
			}

			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			fmt.Println("上滑")
			robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
		case <-chADEnd:
			break loop
		}
	}

	var closeBtnLT, closeBtnRB robotgo.Point
	closeBtnLT.X = ocr.AppWidth*2 - 45
	closeBtnLT.Y = 50
	closeBtnRB.X = ocr.AppWidth*2 - 45 + 23
	closeBtnRB.Y = 50 + 23
	fmt.Printf("点击 关闭(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
	MoveClickTitle(closeBtnLT, closeBtnRB)
	robotgo.Sleep(2)
}

func watchVideoAD() {
	fmt.Println("watchVideoAD")
	predictADEndTick = 0
	go monitor("watchVideoAD", func() bool {
		bCountdownComplete := true
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if strings.Contains(txt, "秒后完成") {
				if predictADEndTick == 0 {
					leftSec, _ := Utils.ExtractNumber(txt)
					predictADEndTick = time.Now().Unix() + int64(leftSec+3)
				}
				bCountdownComplete = false
				break
			}
		}
		return bCountdownComplete && time.Now().Unix() > predictADEndTick
	})
	chADStart <- struct{}{}

	startTick := time.Now()
loop:
	for {
		if time.Now().Unix()-startTick.Unix() > 120 {
			fmt.Println("timeout")
			break loop
		}
		select {
		case <-time.After((10 + time.Duration(Utils.R.Intn(2))) * time.Second):
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			fmt.Println("上滑")
			robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
		case <-chADEnd:
			break loop
		}
	}

	newX := ocr.AppX + 56/2 + Utils.R.Intn(16/2)
	newY := ocr.AppY + 64/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)
}

func watchScroll6sAD() {
	fmt.Println("watchScroll6sAD")
	predictADEndTick = 0
	newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
	newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
	robotgo.Move(newX, newY)
	fmt.Println("上滑")
	robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)

	go monitor("watchScroll6sAD", func() bool {
		bCountdownComplete := true
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if strings.Contains(txt, "滑动浏览") {
				if predictADEndTick == 0 {
					leftSec, _ := Utils.ExtractNumber(txt)
					predictADEndTick = time.Now().Unix() + int64(leftSec+3)
				}
				bCountdownComplete = false
				break
			}
		}
		return bCountdownComplete && time.Now().Unix() > predictADEndTick
	})
	chADStart <- struct{}{}

	startTick := time.Now()
loop:
	for {
		if time.Now().Unix()-startTick.Unix() > 120 {
			fmt.Println("timeout")
			break loop
		}
		select {
		case <-time.After((6 + time.Duration(Utils.R.Intn(2))) * time.Second):
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth*3/4)
			newY := ocr.AppY + ocr.AppHeight*3/4 + Utils.R.Intn(ocr.AppHeight/4)
			robotgo.Move(newX, newY)
			fmt.Println("上滑")
			robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
		case <-chADEnd:
			break loop
		}
	}

	newX = ocr.AppX + 30/2 + Utils.R.Intn(14/2)
	newY = ocr.AppY + 50/2 + Utils.R.Intn(22/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)
}

func watchScroll10sAD() {
	fmt.Println("watchScroll10sAD")
	predictADEndTick = 0
	newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
	newY := ocr.AppY + ocr.AppHeight*3/4 + Utils.R.Intn(ocr.AppHeight/4)
	robotgo.Move(newX, newY)
	fmt.Println("上滑")
	robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)

	go monitor("watchScroll10sAD", func() bool {
		bCountdownComplete := false
		bCountdownTips := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if strings.Contains(txt, "奖励已发放") {
				bCountdownComplete = true
			} else if strings.Contains(txt, "秒后发放奖励") {
				if predictADEndTick == 0 {
					leftSec, _ := Utils.ExtractNumber(txt)
					predictADEndTick = time.Now().Unix() + int64(leftSec+3)
				}
				bCountdownTips = true
			} else if strings.Contains(txt, "计时已暂停，上滑继续") {
				bCountdownTips = true
			}
		}
		return bCountdownComplete && !bCountdownTips && time.Now().Unix() > predictADEndTick
	})
	chADStart <- struct{}{}

	startTick := time.Now()
loopSlideAD:
	for {
		if time.Now().Unix()-startTick.Unix() > 120 {
			fmt.Println("timeout")
			break loopSlideAD
		}

		select {
		case <-time.After((10 + time.Duration(Utils.R.Intn(2))) * time.Second):
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight*3/4 + Utils.R.Intn(ocr.AppHeight/4)
			robotgo.Move(newX, newY)
			fmt.Println("上滑")
			robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
		case <-chADEnd:
			break loopSlideAD
		}
	}

	var closeBtnLT, closeBtnRB robotgo.Point
	closeBtnLT.X = ocr.AppWidth*2 - 50
	closeBtnLT.Y = 48
	closeBtnRB.X = ocr.AppWidth*2 - 50 + 20
	closeBtnRB.Y = 48 + 20
	fmt.Printf("点击 关闭(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
	MoveClickTitle(closeBtnLT, closeBtnRB)
	robotgo.Sleep(2)
}

func watchClickToSkipAD() {
	fmt.Println("watchClickToSkipAD")
	go monitor("watchClickToSkipAD", func() bool {
		bClickToSkip := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if txt == "点击一下领奖励" || txt == "点击按钮可立即领取奖励" {
				bClickToSkip = true
				break
			}
		}
		return bClickToSkip
	})
	chADStart <- struct{}{}
loop1:
	for {
		select {
		case <-time.After((10 + time.Duration(Utils.R.Intn(2))) * time.Second):
			// newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			// newY := ocr.AppY + ocr.AppHeight*3/4 + Utils.R.Intn(ocr.AppHeight/4)
			// robotgo.Move(newX, newY)
			// fmt.Println("上滑")
			// robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
		case <-chADEnd:
			break loop1
		}
	}

	if !OCRMoveClickTitle("点击一下领奖励", 0) {
		OCRMoveClickTitle("点击按钮可立即领取奖励", 0)
	}

	// 此时跳转至其他App
	go monitor("watchInteractiveAD", func() bool {
		bTaoliveBtn := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if txt == "点淘" {
				bTaoliveBtn = true
				break
			}
		}
		return bTaoliveBtn
	})
	chADStart <- struct{}{}
loop2:
	for {
		select {
		case <-time.After((10 + time.Duration(Utils.R.Intn(2))) * time.Second):
			// newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			// newY := ocr.AppY + ocr.AppHeight*3/4 + Utils.R.Intn(ocr.AppHeight/4)
			// robotgo.Move(newX, newY)
			// fmt.Println("上滑")
			// robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
		case <-chADEnd:
			break loop2
		}
	}

	OCRMoveClickTitle("点淘", 0)

	// 此时跳转回点淘广告界面
	go monitor("watchClickToSkipAD", func() bool {
		bADComplete := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if txt == "点击广告可领取奖励|跳过" || txt == "奖励已领取|跳过" {
				bADComplete = true
				break
			}
		}
		return bADComplete
	})
	chADStart <- struct{}{}
loop3:
	for {
		select {
		case <-time.After((10 + time.Duration(Utils.R.Intn(2))) * time.Second):
			// newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			// newY := ocr.AppY + ocr.AppHeight*3/4 + Utils.R.Intn(ocr.AppHeight/4)
			// robotgo.Move(newX, newY)
			// fmt.Println("上滑")
			// robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
		case <-chADEnd:
			break loop3
		}
	}

	var backBtnLT, backeBtnRB robotgo.Point
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		Polygon := v.([]interface{})[0]
		if txt == "点击广告可领取奖励|跳过" || txt == "奖励已领取|跳过" {
			h := int(Polygon.([]interface{})[2].([]interface{})[1].(float64)) - int(Polygon.([]interface{})[1].([]interface{})[1].(float64))
			backBtnLT.X = int(Polygon.([]interface{})[1].([]interface{})[0].(float64)) - 2*h
			backBtnLT.Y = int(Polygon.([]interface{})[1].([]interface{})[1].(float64))
			backeBtnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
			backeBtnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			// fmt.Println(x, y)
			// 点击 跳过
			fmt.Printf("点击 跳过(%3d, %3d)-(%3d, %3d)\n", backBtnLT.X, backBtnLT.Y, backeBtnRB.X, backeBtnRB.Y)
			MoveClickTitle(backBtnLT, backeBtnRB)
			robotgo.Sleep(2)
			break
		}
	}
}

func watchSearchScrollAD() {
	fmt.Println("watchSearchScrollAD")
	// 搜索发现左下角坐标，搜索看更多左上角坐标作为随机点击搜索关键词的区域
	var discoveryLB, searchMoreLT, keyWordLT, keyWordRB robotgo.Point
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		Polygon := v.([]interface{})[0]
		if txt == "搜索发现" {
			discoveryLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
			discoveryLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
			keyWordLT.X = discoveryLB.X
			keyWordLT.Y = discoveryLB.Y + 16
		} else if txt == "搜索看更多" {
			searchMoreLT.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
			searchMoreLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
			keyWordRB.X = searchMoreLT.X - 6
			keyWordRB.Y = keyWordLT.Y + 578
		}
	}

	// 搜索发现下6行两列共12个关键词，设定区域防止点到空白地方，横向间隔10px，纵向间隔14px
	candidateArr := make([]PtPair, 0, 1)
	for i := 0; i < 6; i++ {
		candidateArr = append(candidateArr,
			PtPair{
				LeftTop: robotgo.Point{
					X: keyWordLT.X,
					Y: keyWordLT.Y + (86+14)*i,
				},
				RightBtm: robotgo.Point{
					X: keyWordLT.X + 264,
					Y: keyWordLT.Y + (86+14)*i + 86,
				},
			},
		)
		candidateArr = append(candidateArr,
			PtPair{
				LeftTop: robotgo.Point{
					X: keyWordLT.X + 264 + 10,
					Y: keyWordLT.Y + (86+14)*i,
				},
				RightBtm: robotgo.Point{
					X: keyWordLT.X + 264 + 10 + 192,
					Y: keyWordLT.Y + (86+14)*i + 86,
				},
			},
		)
	}

	i := Utils.R.Intn(len(candidateArr))
	fmt.Printf("点击 随意一个关键词(%3d, %3d)-(%3d, %3d)\n", candidateArr[i].LeftTop.X, candidateArr[i].LeftTop.Y, candidateArr[i].RightBtm.X, candidateArr[i].RightBtm.Y)
	MoveClickTitle(candidateArr[i].LeftTop, candidateArr[i].RightBtm)
	robotgo.Sleep(2)

	watchScroll6sAD()

	newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + 48/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)
}

func watchInteractiveAD() {
	fmt.Println("watchInteractiveAD")
	go monitor("watchInteractiveAD", func() bool {
		bClickToSkip := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			Polygon := v.([]interface{})[0]
			if txt == "免费获取" && int(Polygon.([]interface{})[1].([]interface{})[0].(float64)) < (ocr.AppX+ocr.AppWidth)*2-100 {
				bClickToSkip = true
				break
			}
		}
		return bClickToSkip
	})
	chADStart <- struct{}{}
loop1:
	for {
		select {
		case <-time.After((6 + time.Duration(Utils.R.Intn(2))) * time.Second):
			if ExistText("恭喜获得特权") {
				for _, v := range ocr.OCRResult {
					txt := v.([]interface{})[1].([]interface{})[0].(string)
					Polygon := v.([]interface{})[0]
					var titleRT, closeBtnLT, closeBtnRB robotgo.Point
					if txt == "恭喜获得特权" {
						titleRT.X = int(Polygon.([]interface{})[1].([]interface{})[0].(float64))
						titleRT.Y = int(Polygon.([]interface{})[1].([]interface{})[1].(float64))
						closeBtnLT.X = titleRT.X + 60/2
						closeBtnLT.Y = titleRT.Y - 58/2
						closeBtnRB.X = titleRT.X + 60/2 + 18/2
						closeBtnRB.Y = closeBtnLT.Y + 18/2
						fmt.Printf("点击 关闭(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
						MoveClickTitle(closeBtnLT, closeBtnRB)
						robotgo.Sleep(2)
					}
				}
			}
		case <-chADEnd:
			break loop1
		}
	}

	var closeBtnLT, closeBtnRB robotgo.Point
	closeBtnLT.X = ocr.AppWidth*2 - 36
	closeBtnLT.Y = 56
	closeBtnRB.X = ocr.AppWidth*2 - 36 + 18
	closeBtnRB.Y = 56 + 18
	fmt.Printf("点击 关闭(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
	MoveClickTitle(closeBtnLT, closeBtnRB)
	robotgo.Sleep(2)
}
