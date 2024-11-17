package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
	"regexp"
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

func matchText(reStr string) bool {
	re := regexp.MustCompile(reStr)
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		if re.MatchString(txt) {
			return true
		}
	}

	return false
}

func waitForLeave(scene string) {
	fmt.Println("-> waitForLeave", scene)
	defer fmt.Println("<- waitForLeave", scene)
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

func processInstantBonus(scene string, sceneEntry string) error {
	if scene == "赚步数" || scene == "得体力" || scene == "赚次数" || scene == "赚钱卡" {
		for ExistText("立即领奖") || ExistText("视频福利") {
			if OCRMoveClickTitle("立即领奖", 0) {
				WatchAD(scene, sceneEntry)
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					return err
				}
			} else if OCRMoveClickTitle("视频福利", 0) {
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					return err
				}

				if OCRMoveClickTitle("立即领奖", 0) {
					WatchAD(scene, sceneEntry)
					err := ocr.Ocr(nil, nil, nil, nil)
					if err != nil {
						return err
					}
				}
			}
		}
	} else if scene == "定提醒" {
		for ExistText("立即领奖") || ExistText("视频福利") || ExistText("看5s得") {
			if OCRMoveClickTitle("立即领奖", 0) {
				WatchAD(scene, sceneEntry)
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					return err
				}
			} else if OCRMoveClickTitle("视频福利", 0) {
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					return err
				}

				if OCRMoveClickTitle("立即领奖", 0) {
					WatchAD(scene, sceneEntry)
					err := ocr.Ocr(nil, nil, nil, nil)
					if err != nil {
						return err
					}
				}
			} else if OCRMoveClickTitle("看5s得", 0) {
				WatchAD(scene, sceneEntry)
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func waitForEnter(scene string, sceneEntry string) {
	fmt.Println("-> waitForEnter", scene)
	defer fmt.Println("<- waitForEnter", scene)
	for {
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if ExistText("继续做任务") {
			OCRMoveClickTitle("继续做任务", 0)
		}

		err = processInstantBonus(scene, sceneEntry)
		if err != nil {
			panic(err)
		}

		if ExistText(scene) {
			break
		} else if sceneEntry != "" && ExistText(sceneEntry) {
			OCRMoveClickTitle(sceneEntry, 0)
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

func WatchAD(lastScene string, sceneEntry string) {
	fmt.Println("-> WatchAD", lastScene)
	defer fmt.Println("<- WatchAD", lastScene)
	chADStart = make(chan struct{})
	chADEnd = make(chan struct{})
	defer close(chADStart)
	defer close(chADEnd)

	if lastScene != "" {
		fmt.Println("等待离开", lastScene)
		waitForLeave(lastScene)
	}

checkAgain:
	if containText("秒后完成") || containText("秒后领奖") || containText("秒后发奖") ||
		containText("后完成") || containText("后领奖") || containText("后发奖") ||
		(containText("秒后发放") && !containText("秒后发放奖励")) {
		watchLiveOrVideoAD(60)
	} else if containText("滑动浏览") {
		watchScroll6sAD()
	} else if containText("秒后发放奖励") || containText("计时已暂停，上滑继续") {
		watchScroll10sAD()
	} else if containText("点击广告可领取奖励|跳过") || containText("点击按钮可立即领取奖励") || containText("点击一下领奖励") {
		watchClickToSkipAD()
	} else if containText("查看详情立即领奖") {
		watchClickInAppToSkipAD()
	} else if (containText("搜索领元宝") || containText("搜索领体力") || containText("搜索领次数") || containText("搜索领步数") || containText("搜索有福利")) && !containText("搜索并点击") {
		watchSearchScrollAD()
	} else if containText("搜索并点击") && containText("个宝贝") {
		watchSearchAndClickAD()
	} else if containText("立即获取") && containText("跳过") {
		watchInteractiveAD()
	} else if matchText(`^\d+s后可领取奖.*$`) || containText("该视频提到的内容是") {
		watchChooseAnswerAD()
	} else if containText("s|跳过") {

	} else {
		fmt.Println("Unknown, 上滑")
		newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
		newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
		robotgo.Move(newX, newY)
		robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
		robotgo.Sleep(2)
		err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		goto checkAgain
		// watchLiveOrVideoAD()
	}

	if lastScene != "" {
		fmt.Println("等待回到", lastScene)
		waitForEnter(lastScene, sceneEntry)
	}
}

func processLivePopup() {
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
		} else if txt == "关注领取" {
			btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
			btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
			btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
			btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			centerX := btnLB.X + int((btnRB.X-btnLB.X)/2)
			closeBtnLT.X = centerX - 22/2
			closeBtnRB.X = centerX + 22/2
			closeBtnLT.Y = btnLB.Y + (1042 - 956)
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
		} else if txt == "关注主播参与" {
			btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
			btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
			btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
			btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			centerX := btnLB.X + int((btnRB.X-btnLB.X)/2)
			closeBtnLT.X = centerX - 22/2
			closeBtnRB.X = centerX + 22/2
			closeBtnLT.Y = btnLB.Y + (894 - 812)
			closeBtnRB.Y = closeBtnLT.Y + 22
			fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
			MoveClickTitle(closeBtnLT, closeBtnRB)
			robotgo.Sleep(2)
			// } else if txt == "双11狂欢节" {
			// 	btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
			// 	btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
			// 	btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
			// 	btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			// 	centerX := btnLB.X + int((btnRB.X-btnLB.X)/2)
			// 	closeBtnLT.X = centerX - 22/2
			// 	closeBtnRB.X = centerX + 22/2
			// 	closeBtnLT.Y = btnLB.Y - (608 - 494)
			// 	closeBtnRB.Y = closeBtnLT.Y + 22
			// 	fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
			// 	MoveClickTitle(closeBtnLT, closeBtnRB)
			// 	robotgo.Sleep(2)
		} else if txt == "一键领取" {
			btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
			btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
			btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
			btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			centerX := btnLB.X + int((btnRB.X-btnLB.X)/2)
			closeBtnLT.X = centerX - 22/2
			closeBtnRB.X = centerX + 22/2
			closeBtnLT.Y = btnLB.Y + (1040 - 940)
			closeBtnRB.Y = closeBtnLT.Y + 22
			fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
			MoveClickTitle(closeBtnLT, closeBtnRB)
			robotgo.Sleep(2)
		}
	}
}

func watchLiveOrVideoAD(duration int64) {
	fmt.Println("-> watchLiveOrVideoAD")
	defer fmt.Println("<- watchLiveOrVideoAD")
retry:
	predictADEndTick = time.Now().Unix() + duration
	isLive := false
	go monitor("watchLiveOrVideoAD", func() bool {
		bCountdownComplete := true
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			re := regexp.MustCompile(`^\d+秒后.*`)
			if re.MatchString(txt) {
				leftSec, _ := Utils.ExtractNumber(txt)
				if predictADEndTick > 0 && predictADEndTick-(ocr.OCRTick+int64(leftSec)) > 5 {
					predictADEndTick = ocr.OCRTick + int64(leftSec)
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
			if ExistText("点击翻倍") {
				OCRMoveClickTitle("点击翻倍", 0)
			}
			if ExistText("点击x4倍") {
				OCRMoveClickTitle("点击x4倍", 0)
			}

			processLivePopup()

			if !isLive && (ExistText("说点什么") || ExistText("说点") || ExistText("关注") || ExistText("评论已关闭")) {
				isLive = true
			}

			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			fmt.Printf("上滑 直播: %v\n", isLive)
			robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
			robotgo.Sleep(2)
		case <-chADEnd:
			break loop
		}
	}

	if isLive {
		for {
			var closeBtnLT, closeBtnRB robotgo.Point
			closeBtnLT.X = ocr.AppWidth*2 - 45
			closeBtnLT.Y = 50
			closeBtnRB.X = ocr.AppWidth*2 - 45 + 23
			closeBtnRB.Y = 50 + 23
			fmt.Printf("点击 关闭(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
			MoveClickTitle(closeBtnLT, closeBtnRB)
			robotgo.Sleep(2)

			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				panic(err)
			}

			if !containText("说点") && !ExistText("关注") && !ExistText("最爱") {
				break
			}

			processLivePopup()
		}
	} else {
		if ExistText("可领奖|关闭广告") {
			var closeBtnLT, closeBtnRB robotgo.Point
			for _, v := range ocr.OCRResult {
				txt := v.([]interface{})[1].([]interface{})[0].(string)
				Polygon := v.([]interface{})[0]
				if strings.Contains(txt, "可领奖|关闭广告") {
					fontHeight := int(Polygon.([]interface{})[2].([]interface{})[1].(float64)) - int(Polygon.([]interface{})[1].([]interface{})[1].(float64))
					closeBtnLT.X = int(Polygon.([]interface{})[1].([]interface{})[0].(float64)) - 2*fontHeight
					closeBtnLT.Y = int(Polygon.([]interface{})[1].([]interface{})[1].(float64))
					closeBtnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
					closeBtnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
					// fmt.Println(x, y)
					// 点击 跳过
					fmt.Printf("点击 关闭广告(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
					MoveClickTitle(closeBtnLT, closeBtnRB)
					robotgo.Sleep(2)
					break
				}
			}
		} else {
			newX := ocr.AppX + 56/2 + Utils.R.Intn(16/2)
			newY := ocr.AppY + 64/2 + Utils.R.Intn(26/2)
			fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
			robotgo.MoveClick(newX, newY)
			robotgo.Sleep(2)
		}
	}

	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}

	if ExistText("继续做任务") {
		OCRMoveClickTitle("继续做任务", 0)
		goto retry
	}
}

// func watchVideoAD() {
// 	fmt.Println("watchVideoAD")
// 	predictADEndTick = 0
// 	go monitor("watchVideoAD", func() bool {
// 		bCountdownComplete := true
// 		for _, v := range ocr.OCRResult {
// 			txt := v.([]interface{})[1].([]interface{})[0].(string)
// 			if strings.Contains(txt, "秒后完成") {
// 				if predictADEndTick == 0 {
// 					leftSec, _ := Utils.ExtractNumber(txt)
// 					predictADEndTick = time.Now().Unix() + int64(leftSec+3)
// 				}
// 				bCountdownComplete = false
// 				break
// 			}
// 		}
// 		return bCountdownComplete && time.Now().Unix() > predictADEndTick
// 	})
// 	chADStart <- struct{}{}

// 	startTick := time.Now()
// loop:
// 	for {
// 		if time.Now().Unix()-startTick.Unix() > 120 {
// 			fmt.Println("timeout")
// 			break loop
// 		}
// 		select {
// 		case <-time.After((10 + time.Duration(Utils.R.Intn(2))) * time.Second):
// 			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
// 			newY := ocr.AppY + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
// 			robotgo.Move(newX, newY)
// 			fmt.Println("上滑")
// 			robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
// 		case <-chADEnd:
// 			break loop
// 		}
// 	}

// 	newX := ocr.AppX + 56/2 + Utils.R.Intn(16/2)
// 	newY := ocr.AppY + 64/2 + Utils.R.Intn(26/2)
// 	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
// 	robotgo.MoveClick(newX, newY)
// 	robotgo.Sleep(2)
// }

func watchScroll6sAD() {
	fmt.Println("watchScroll6sAD")
retry:
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
			re := regexp.MustCompile(`^滑动浏览\d+.*`)
			if re.MatchString(txt) {
				leftSec, _ := Utils.ExtractNumber(txt)
				if predictADEndTick == 0 && leftSec <= 60 {
					predictADEndTick = ocr.OCRTick + int64(leftSec)
				}
				bCountdownComplete = false
				break
			}
		}
		fmt.Printf("Now: %d, predict: %d\n", time.Now().Unix(), predictADEndTick)
		return bCountdownComplete && predictADEndTick > 0 && time.Now().Unix() > predictADEndTick
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

	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}

	if ExistText("继续做任务") {
		OCRMoveClickTitle("继续做任务", 0)
		goto retry
	}
}

func watchScroll10sAD() {
	fmt.Println("watchScroll10sAD")
retry:
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
				re := regexp.MustCompile(`^\d+秒后发放奖励$`)
				if re.MatchString(txt) {
					leftSec, _ := Utils.ExtractNumber(txt)
					if predictADEndTick == 0 && leftSec <= 30 {
						predictADEndTick = ocr.OCRTick + int64(leftSec)
					}
					bCountdownComplete = false
					break
				}
				bCountdownTips = true
			} else if strings.Contains(txt, "计时已暂停，上滑继续") {
				bCountdownTips = true
			}
		}
		return bCountdownComplete && !bCountdownTips && predictADEndTick > 0 && time.Now().Unix() > predictADEndTick
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
			for _, v := range ocr.OCRResult {
				txt := v.([]interface{})[1].([]interface{})[0].(string)
				Polygon := v.([]interface{})[0]
				var btnLB, btnRB, closeBtnLT, closeBtnRB robotgo.Point
				if txt == "立即抽奖" {
					btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
					btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
					btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
					btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
					closeBtnLT.X = btnRB.X + (1198 - 1068)
					closeBtnRB.X = closeBtnLT.X + 26
					closeBtnLT.Y = btnRB.Y - (920 - 254)
					closeBtnRB.Y = closeBtnLT.Y + 26
					fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
					MoveClickTitle(closeBtnLT, closeBtnRB)
					robotgo.Sleep(2)
					break
				}
			}

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

	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}

	if ExistText("继续做任务") {
		OCRMoveClickTitle("继续做任务", 0)
		goto retry
	}
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

	if !OCRMoveClickTitle("查看详情", 0) {
		OCRMoveClickTitle("取消", 0)
	}

	// 此时跳转至其他App
	go monitor("watchInteractiveAD", func() bool {
		bTaoliveBtn := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if strings.Contains(txt, "点淘") || strings.Contains(txt, "完成") {
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

	if !OCRMoveClickTitle("点淘", 0) {
		OCRMoveClickTitle("完成", 0)
	}

	// 此时跳转回点淘广告界面
	go monitor("watchClickToSkipAD", func() bool {
		bADComplete := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if strings.Contains(txt, "点击广告可领取奖励|跳过") || strings.Contains(txt, "奖励已领取|跳过") || strings.Contains(txt, "点击按钮可立即领取奖励") {
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

	var backBtnLT, backBtnRB robotgo.Point
	for _, v := range ocr.OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		Polygon := v.([]interface{})[0]
		if strings.Contains(txt, "点击广告可领取奖励|跳过") || strings.Contains(txt, "奖励已领取|跳过") {
			fontHeight := int(Polygon.([]interface{})[2].([]interface{})[1].(float64)) - int(Polygon.([]interface{})[1].([]interface{})[1].(float64))
			backBtnLT.X = int(Polygon.([]interface{})[1].([]interface{})[0].(float64)) - 2*fontHeight
			backBtnLT.Y = int(Polygon.([]interface{})[1].([]interface{})[1].(float64))
			backBtnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
			backBtnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
			// fmt.Println(x, y)
			// 点击 跳过
			fmt.Printf("点击 跳过(%3d, %3d)-(%3d, %3d)\n", backBtnLT.X, backBtnLT.Y, backBtnRB.X, backBtnRB.Y)
			MoveClickTitle(backBtnLT, backBtnRB)
			robotgo.Sleep(2)
			break
		}
	}
}

func watchClickInAppToSkipAD() {
	fmt.Println("watchClickInAppToSkipAD")
	go monitor("watchClickInAppToSkipAD", func() bool {
		bClickToSkip := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if txt == "查看详情立即领奖" {
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

	OCRMoveClickTitle("查看详情", 0)

	// 此时跳转至点淘App内商品详情
	newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + 48/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)

	// 此时跳转回点淘广告界面
	go monitor("watchClickInAppToSkipAD", func() bool {
		bADComplete := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if txt == "奖励已发放" || txt == "查看详情立即领奖" {
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

	var closeBtnLT, closeBtnRB robotgo.Point
	closeBtnLT.X = ocr.AppWidth*2 - 56
	closeBtnLT.Y = 56
	closeBtnRB.X = ocr.AppWidth*2 - 56 + 18
	closeBtnRB.Y = 56 + 18
	fmt.Printf("点击 关闭(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
	MoveClickTitle(closeBtnLT, closeBtnRB)
	robotgo.Sleep(2)
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

func watchSearchAndClickAD() {
	fmt.Println("watchSearchAndClickAD")
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

retry:
	go monitor("watchSearchAndClickAD", func() bool {
		bClickComplete := true
		if !ExistText("搜索领元宝") {
			return false
		}
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if strings.Contains(txt, "点击3个") {
				bClickComplete = false
				break
			}
		}
		return bClickComplete
	})
	chADStart <- struct{}{}

	startTick := time.Now()
loop:
	for {
		if time.Now().Unix()-startTick.Unix() > 60 {
			fmt.Println("timeout")
			break loop
		}
		select {
		case <-time.After(time.Duration(Utils.R.Intn(2)) * time.Second):
			inMerchantDetail := false
			for _, v := range ocr.OCRResult {
				txt := v.([]interface{})[1].([]interface{})[0].(string)
				Polygon := v.([]interface{})[0]
				var detailBtnLT, detailBtnRB robotgo.Point
				if strings.Contains(txt, "下单约得") {
					detailBtnLT.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
					detailBtnLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
					detailBtnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
					detailBtnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
					fmt.Printf("点击 下单约得(%3d, %3d)-(%3d, %3d)\n", detailBtnLT.X, detailBtnLT.Y, detailBtnRB.X, detailBtnRB.Y)
					MoveClickTitle(detailBtnLT, detailBtnRB)
					robotgo.Sleep(2)
					inMerchantDetail = true
					break
				}
			}

			if inMerchantDetail {
				// 此时跳转至点淘App内商品详情
				newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
				newY := ocr.AppY + 48/2 + Utils.R.Intn(26/2)
				fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
				robotgo.MoveClick(newX, newY)
				robotgo.Sleep(2)
			}

			// newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth*3/4)
			// newY := ocr.AppY + ocr.AppHeight*3/4 + Utils.R.Intn(ocr.AppHeight/4)
			// robotgo.Move(newX, newY)
			// fmt.Println("上滑")
			// robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
		case <-chADEnd:
			break loop
		}
	}

	newX := ocr.AppX + 28/2 + Utils.R.Intn(14/2)
	newY := ocr.AppY + 48/2 + Utils.R.Intn(26/2)
	fmt.Printf("点击 返回(%3d, %3d)\n", newX, newY)
	robotgo.MoveClick(newX, newY)
	robotgo.Sleep(2)

	err := ocr.Ocr(nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}

	if ExistText("继续做任务") {
		OCRMoveClickTitle("继续做任务", 0)
		goto retry
	}
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
						closeBtnLT.X = titleRT.X + 60
						closeBtnLT.Y = titleRT.Y - 58
						closeBtnRB.X = titleRT.X + 60 + 18
						closeBtnRB.Y = closeBtnLT.Y + 18
						fmt.Printf("点击 恭喜获得特权 关闭(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
						MoveClickTitle(closeBtnLT, closeBtnRB)
						robotgo.Sleep(2)
					}
				}
			}
		case <-chADEnd:
			break loop1
		}
	}

	// var closeBtnLT, closeBtnRB robotgo.Point
	// for _, v := range ocr.OCRResult {
	// 	txt := v.([]interface{})[1].([]interface{})[0].(string)
	// 	Polygon := v.([]interface{})[0]
	// 	if txt == "免费获取" && int(Polygon.([]interface{})[1].([]interface{})[0].(float64)) < (ocr.AppX+ocr.AppWidth)*2-100 {
	// 		closeBtnLT.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
	// 		closeBtnLT.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
	// 		closeBtnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
	// 		closeBtnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
	// 		break
	// 	}
	// }

	var closeBtnLT, closeBtnRB robotgo.Point
	closeBtnLT.X = ocr.AppWidth*2 - 54
	closeBtnLT.Y = 56
	closeBtnRB.X = ocr.AppWidth*2 - 54 + 18
	closeBtnRB.Y = 56 + 18
	fmt.Printf("点击 关闭(%3d, %3d)-(%3d, %3d)\n", closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
	MoveClickTitle(closeBtnLT, closeBtnRB)
	robotgo.Sleep(2)
}

func watchChooseAnswerAD() {
	fmt.Println("watchChooseAnswerAD")
	go monitor("watchChooseAnswerAD", func() bool {
		bCountdownComplete := false
		for _, v := range ocr.OCRResult {
			txt := v.([]interface{})[1].([]interface{})[0].(string)
			if txt == "已领取奖励" {
				bCountdownComplete = true
				break
			}
		}
		return bCountdownComplete
	})
	chADStart <- struct{}{}
loop1:
	for {
		select {
		case <-time.After(time.Duration(Utils.R.Intn(2)) * time.Second):
			if ExistText("喜马拉雅") && (ExistText("立春") || ExistText("春分") || ExistText("立夏") || ExistText("夏至") || ExistText("立秋") || ExistText("秋分") ||
				ExistText("立冬") || ExistText("冬至") || ExistText("小暑") || ExistText("大暑") || ExistText("处暑") || ExistText("小寒") ||
				ExistText("大寒") || ExistText("雨水") || ExistText("谷雨") || ExistText("白露") || ExistText("寒露") || ExistText("霜降") ||
				ExistText("小雪") || ExistText("大雪") || ExistText("惊蛰") || ExistText("清明") || ExistText("小满") || ExistText("芒种")) {
				OCRMoveClickTitle("喜马拉雅", 0)
			}
		case <-chADEnd:
			break loop1
		}
	}

	OCRMoveClickTitle("跳过", 0)
}
