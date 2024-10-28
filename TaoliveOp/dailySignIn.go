package TaoliveOp

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/Utils"
	"fmt"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

func DoDailySignIn() error {
	bNeedScroll := true
	for bNeedScroll {
		tasks, err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for _, v := range tasks {
			txt := v.([]interface{})[1].([]interface{})[0]
			if txt == "去完成" {
				bNeedScroll = false
				break
			}
		}

		if bNeedScroll {
			// 从下往上滑动
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + 38 + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(10) + 90), 3, 50, Utils.R.Intn(10)-5)
			robotgo.Sleep(1)
		}
	}

	allDone := false
	for !allDone {
		tasks, err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		for _, v := range tasks {
			txt := v.([]interface{})[1].([]interface{})[0]
			if txt == "去完成" {
				Polygon := v.([]interface{})[0]
				// fmt.Println(Polygon.([]interface{})[0].([]interface{})[0].(float64))
				var leftTop, rightBtm robotgo.Point
				leftTop.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
				leftTop.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
				rightBtm.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
				rightBtm.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
				// fmt.Println(x, y)
				// 点击 去完成
				fmt.Println("点击 去完成")
				MoveClickTitle(leftTop, rightBtm, 0)
				robotgo.Sleep(1)
				break
			}
		}

		Ads, err := ocr.Ocr(nil, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if ADType(Ads) == ADType_Live {
			// 关闭按钮 的中心点在 更多直播 最后一个字的右上角坐标往上54pixel，宽高22px
			adTime := 60
			var closeBtnLT, closeBtnRM robotgo.Point
			for _, v := range Ads {
				txt := v.([]interface{})[1].([]interface{})[0].(string)
				if strings.Contains(txt, "秒后完成") {
					adTime, err = Utils.ExtractNumber(txt)
					if err != nil {
						panic(err)
					}
				} else if txt == "更多直播" {
					Polygon := v.([]interface{})[0]
					closeBtnLT.X = int(Polygon.([]interface{})[1].([]interface{})[0].(float64) - 22/2)
					closeBtnLT.Y = int(Polygon.([]interface{})[1].([]interface{})[1].(float64) - 54 - 22/2)
					closeBtnRM.X = int(Polygon.([]interface{})[1].([]interface{})[0].(float64) + 22/2)
					closeBtnRM.Y = int(Polygon.([]interface{})[1].([]interface{})[1].(float64) - 54 + 22/2)
				}
			}
			fmt.Printf("等待 %d 秒后关闭广告\n", adTime)
			robotgo.Sleep(adTime)
			MoveClickTitle(closeBtnLT, closeBtnRM, 0)
		} else if ADType(Ads) == ADType_Scroll6s {
			adTime := 60
			for _, v := range Ads {
				txt := v.([]interface{})[1].([]interface{})[0].(string)
				if strings.Contains(txt, "滑动浏览") {
					adTime, err = Utils.ExtractNumber(txt)
					if err != nil {
						panic(err)
					}
					break
				}
			}

			adEnd := make(chan struct{})
			go func() {
				for {
					Ads, err := ocr.Ocr(nil, nil, nil, nil)
					if err != nil {
						panic(err)
					}

					needScroll := false
					for _, v := range Ads {
						txt := v.([]interface{})[1].([]interface{})[0].(string)
						if strings.Contains(txt, "滑动浏览") {
							fmt.Println(txt)
							needScroll = true
							break
						}
					}

					if !needScroll {
						break
					}
				}

				fmt.Println("广告结束")
				adEnd <- struct{}{}
			}()

			fmt.Printf("上滑 并等待 %d 秒后关闭广告\n", adTime)
			newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
			newY := ocr.AppY + 38 + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
			robotgo.Move(newX, newY)
			robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)

		loop:
			for {
				select {
				case <-time.After(6 * time.Second):
					newX := ocr.AppX + Utils.R.Intn(ocr.AppWidth)
					newY := ocr.AppY + 38 + ocr.AppHeight/2 + Utils.R.Intn(ocr.AppHeight/2)
					robotgo.Move(newX, newY)
					fmt.Println("上滑")
					robotgo.ScrollSmooth(-(Utils.R.Intn(30) + 150), 3, 50, Utils.R.Intn(10)-5)
				case <-adEnd:
					break loop
				}
			}

			close(adEnd)
			newX = ocr.AppX + 10 + 30/2 + Utils.R.Intn(14/2)
			newY = ocr.AppY + 38 + 50/2 + Utils.R.Intn(22/2)
			robotgo.MoveClick(newX, newY)
			robotgo.Sleep(1)
		}
	}

	return nil
}
