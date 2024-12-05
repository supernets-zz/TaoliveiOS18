// package main

// import (
// 	"time"

// 	"github.com/go-vgo/robotgo"
// )

// func main() {
// 	robotgo.MouseSleep = 100

// 	time.Sleep(5 * time.Second)
// 	robotgo.ScrollSmooth(-300, 6, 50, 0)
// 	// robotgo.ScrollDir(1000, "down")
// 	// robotgo.ScrollDir(20, "right")

// 	robotgo.Scroll(0, -10)
// 	// robotgo.Scroll(100, 0)

// 	// robotgo.MilliSleep(100)
// 	// robotgo.ScrollSmooth(-10, 6)
// 	// // robotgo.ScrollRelative(10, -100)

// 	// robotgo.Move(10, 20)
// 	// robotgo.MoveRelative(0, -10)
// 	// robotgo.DragSmooth(10, 10)

// 	// robotgo.Click("wheelRight")
// 	// robotgo.Click("left", true)
// 	// robotgo.MoveSmooth(100, 200, 1.0, 10.0)

//		// robotgo.Toggle("left")
//		// robotgo.Toggle("left", "up")
//	}
// package main

// import (
// 	"fmt"
// 	"strconv"

// 	"github.com/go-vgo/robotgo"
// 	"github.com/vcaesar/imgo"
// )

// func main() {
// 	x, y := robotgo.Location()
// 	fmt.Println("pos: ", x, y)

// 	color := robotgo.GetPixelColor(100, 200)
// 	fmt.Println("color---- ", color)

// 	sx, sy := robotgo.GetScreenSize()
// 	fmt.Println("get screen size: ", sx, sy)

// 	bit := robotgo.CaptureScreen(10, 10, 30, 30)
// 	defer robotgo.FreeBitmap(bit)

// 	img := robotgo.ToImage(bit)
// 	imgo.Save("test.png", img)

// 	num := robotgo.DisplaysNum()
// 	fmt.Println("DisplaysNum: ", num)
// 	for i := 0; i < num; i++ {
// 		robotgo.DisplayID = i
// 		img1, err := robotgo.CaptureImg()
// 		if err != nil {
// 			panic(err)
// 		}
// 		path1 := "save_" + strconv.Itoa(i)
// 		robotgo.Save(img1, path1+".png")
// 		robotgo.SaveJpeg(img1, path1+".jpeg", 50)

// 		img2, _ := robotgo.CaptureImg(10, 10, 20, 20)
// 		robotgo.Save(img2, "test_"+strconv.Itoa(i)+".png")

// 		x, y, w, h := robotgo.GetDisplayBounds(i)
// 		fmt.Printf("x: %d, y: %d, w: %d, h: %d\n", x, y, w, h)
// 		img3, err := robotgo.CaptureImg(x, y, w, h)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println("Capture error: ", err)
// 		robotgo.Save(img3, path1+"_1.png")
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"image"

// 	"github.com/andybrewer/mack"
// 	"github.com/go-vgo/robotgo"
// )

// func otsu(img *image.Gray) uint32 {
// 	var hist [256]int
// 	bounds := img.Bounds()
// 	for x := bounds.Min.X; x < bounds.Max.X; x++ {
// 		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
// 			hist[img.GrayAt(x, y).Y]++
// 		}
// 	}

// 	total := bounds.Max.X * bounds.Max.Y
// 	var sum float64
// 	for i := 0; i < 256; i++ {
// 		sum += float64(i) * float64(hist[i])
// 	}
// 	var sumB float64
// 	wB := 0
// 	wF := 0
// 	var varMax float64
// 	threshold := 0

// 	for t := 0; t < 256; t++ {
// 		wB += hist[t]
// 		if wB == 0 {
// 			continue
// 		}
// 		wF = total - wB
// 		if wF == 0 {
// 			break
// 		}
// 		sumB += float64(t) * float64(hist[t])

// 		mB := sumB / float64(wB)
// 		mF := (sum - sumB) / float64(wF)
// 		var between float64 = float64(wB) * float64(wF) * (mB - mF) * (mB - mF)
// 		if between >= varMax {
// 			threshold = t
// 			varMax = between
// 		}
// 	}

// 	return uint32(threshold)
// }

// func main() {
// 	// names, err := robotgo.FindNames()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// fmt.Println(names)
// 	fmt.Println(robotgo.GetScaleSize(0))
// 	fpid, err := robotgo.FindIds("iPhone Mirroring")
// 	if err == nil {
// 		fmt.Println("pids... ", fpid)

// 		// if len(fpid) > 0 {
// 		// 	robotgo.TypeStr("Hi galaxy!", fpid[0])
// 		// 	robotgo.KeyTap("a", fpid[0], "cmd")

// 		// 	robotgo.KeyToggle("a", fpid[0])
// 		// 	robotgo.KeyToggle("a", fpid[0], "up")

// 		robotgo.ActivePid(fpid[0])
// 		fmt.Println(robotgo.GetDisplayBounds(0))
// 		response, err := mack.Tell("System Events", fmt.Sprintf("set _P to a reference to (processes whose unix id is %d)", fpid[0]), "set _W to a reference to windows of _P", "[_P's name, _W's size, _W's position]")
// 		if err != nil {
// 			fmt.Println("couldnt get dimensions and position")
// 		} else {
// 			fmt.Println(response)
// 			var appName1, appName2 string
// 			var w, h, x, y int
// 			n, err := fmt.Sscanf(response, "%s %s %d, %d, %d, %d", &appName1, &appName2, &w, &h, &x, &y)
// 			if err != nil {
// 				panic(err)
// 			}
// 			fmt.Println(n, x, y, w, h)
// 			// robotgo.Move(x+int(w/2), y+int(h/2))
// 			img, err := robotgo.CaptureImg(x, y+40, w, h-40)
// 			if err != nil {
// 				panic(err)
// 			}

// 			err = robotgo.SavePng(img, "./t.png")
// 			if err != nil {
// 				panic(err)
// 			}

// 			// // 创建一个新的 Tesseract 客户端
// 			// client := gosseract.NewClient()
// 			// defer client.Close()

// 			// // 设置要识别的图片路径
// 			// // fmt.Println(len(robotgo.ToByteImg(img, "jpeg")))
// 			// // err = client.SetImageFromBytes(robotgo.ToByteImg(img, "jpeg"))
// 			// err = client.SetImage("./screenshot-20241025-144750.png")
// 			// if err != nil {
// 			// 	panic(err)
// 			// }

// 			// // // 获取识别结果，包括文本和位置信息
// 			// // text, err := client.Text() // 得到识别的文本
// 			// // // _, err = client.Text()
// 			// // if err != nil {
// 			// // 	panic(err)
// 			// // }

// 			// // 获取位置信息
// 			// boxes, err := client.GetBoundingBoxes(gosseract.RIL_SYMBOL) // 得到字符的位置信息
// 			// if err != nil {
// 			// 	panic(err)
// 			// }

// 			// // // 输出识别的文本
// 			// // fmt.Println("Recognized text:", text)

// 			// // 输出字符位置
// 			// fmt.Println("Character boxes:")
// 			// for _, box := range boxes {
// 			// 	fmt.Printf("Character: %s, Position: (%d, %d), Width: %d, Height: %d\n",
// 			// 		box.Word, box.Box.Min.X, box.Box.Min.X, box.Box.Dx(), box.Box.Dy())
// 			// }

// 			// gray := image.NewGray(img.Bounds())
// 			// for x := gray.Bounds().Min.X; x < gray.Bounds().Max.X; x++ {
// 			// 	for y := gray.Bounds().Min.Y; y < gray.Bounds().Max.Y; y++ {
// 			// 		r, g, b, _ := img.At(x, y).RGBA()
// 			// 		grayColor := color.Gray{uint8(r+g+b) / 3}
// 			// 		gray.Set(x, y, grayColor)
// 			// 	}
// 			// }

// 			// //分割图片
// 			// bounds := gray.Bounds()
// 			// threshold := otsu(gray) // OTSU算法获取阈值
// 			// binary := image.NewGray(bounds)
// 			// for x := bounds.Min.X; x < bounds.Max.X; x++ {
// 			// 	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
// 			// 		if uint32(gray.GrayAt(x, y).Y) > threshold {
// 			// 			binary.Set(x, y, color.Gray{255})
// 			// 		} else {
// 			// 			binary.Set(x, y, color.Gray{0})
// 			// 		}
// 			// 	}
// 			// }

// 			// client := gosseract.NewClient()
// 			// defer client.Close()

// 			// texts := make([]string, 0)
// 			// bounds := binary.Bounds()
// 			// for x := bounds.Min.X; x < bounds.Max.X; x++ {
// 			// 	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
// 			// 		if binary.GrayAt(x, y).Y == 255 {
// 			// 			continue
// 			// 		}
// 			// 		sx := X
// 			// 		sy := y
// 			// 		ex := x
// 			// 		ey := y
// 			// 		for ; ex < bounds.Max.X && binary.GrayAt(ex, y).Y == 0; ex++ {
// 			// 		}
// 			// 		for ; ey < bounds.Max.Y && binary.GrayAt(x, ey).Y == 0; ey++ {
// 			// 		}
// 			// 		rect := image.Rect(sx, sy, ex, ey)
// 			// 		subimg := binary.SubImage(rect)

// 			// 		pix := subImg.Bounds().Max.X * subImg.Bounds().Max.Y
// 			// 		blackNum := 0
// 			// 		for i := subImg.Bounds().Min.X; i < subImg.Bounds().Max.X; i++ {
// 			// 			for j := subImg.Bounds().Min.Y; j < subImg.Bounds().Max.Y; j++ {
// 			// 				if subImg.At(i, j) == color.Gray{255} {
// 			// 					blackNum++
// 			// 				}
// 			// 			}
// 			// 		}
// 			// 		if float64(blackNum)/float64(pix) < 0.1 { // 去除噪音
// 			// 			continue
// 			// 		}

// 			// 		output, _ := client.ImageToText(subImg)
// 			// 		output = strings.ReplaceAll(output, "\n", "")
// 			// 		output = strings.ReplaceAll(output, " ", "")
// 			// 		texts = append(texts, output)
// 			// 	}
// 			// }

// 			// fmt.Println(texts)
// 		}
// 		// 	robotgo.Kill(fpid[0])
// 		// }
// 	}

// 	// robotgo.ActiveName("chrome")

// 	// isExist, err := robotgo.PidExists(100)
// 	// if err == nil && isExist {
// 	// 	fmt.Println("pid exists is", isExist)

// 	// 	robotgo.Kill(100)
// 	// }

// 	// abool := robotgo.Alert("test", "robotgo")
// 	// if abool {
// 	// 	fmt.Println("ok@@@ ", "ok")
// 	// }

// 	// title := robotgo.GetTitle()
// 	// fmt.Println("title@@@ ", title)
// }

package main

import (
	ocr "TaoliveiOS18/OCR"
	"TaoliveiOS18/TaoliveOp"
	"TaoliveiOS18/Utils"
	IPhoneOp "TaoliveiOS18/iPhoneOp"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-vgo/robotgo"
)

type ResultItem struct {
	Polygon           [][2]float64
	TextAndConfidence []interface{}
}

func main() {
	Utils.R = rand.New(rand.NewSource(time.Now().UnixNano()))
	robotgo.MouseSleep = 100
	robotgo.KeySleep = 100

	if err := IPhoneOp.GetIPhoneMirroringGeometry(); err != nil {
		panic(err)
	}

	// var closeBtnLT, closeBtnRM robotgo.Point
	// closeBtnLT.X = ocr.AppX + ocr.AppWidth - 30/2 - 20/2
	// closeBtnLT.Y = ocr.AppY + 48/2
	// closeBtnRM.X = ocr.AppX + ocr.AppWidth - 30/2
	// closeBtnRM.Y = ocr.AppY + 48/2 + 20/2
	// robotgo.Move(closeBtnRM.X, closeBtnRM.Y)
	// return

	// err := ocr.Ocr(nil, nil, nil, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// TaoliveOp.WatchAD("")
	// return

	// err := ocr.Ocr(nil, nil, nil, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, v := range ocr.OCRResult {
	// 	txt := v.([]interface{})[1].([]interface{})[0].(string)
	// 	Polygon := v.([]interface{})[0]
	// 	var btnLB, btnRB, closeBtnLT, closeBtnRB robotgo.Point
	// 	if txt == "立即抽奖" {
	// 		btnLB.X = int(Polygon.([]interface{})[3].([]interface{})[0].(float64))
	// 		btnLB.Y = int(Polygon.([]interface{})[3].([]interface{})[1].(float64))
	// 		btnRB.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
	// 		btnRB.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
	// 		closeBtnLT.X = btnRB.X + 1198 - 1068
	// 		closeBtnRB.X = closeBtnLT.X + 26
	// 		closeBtnLT.Y = btnRB.Y - (920 - 254)
	// 		closeBtnRB.Y = closeBtnLT.Y + 26
	// 		fmt.Printf("点击 %s(%3d, %3d)-(%3d, %3d)\n", txt, closeBtnLT.X, closeBtnLT.Y, closeBtnRB.X, closeBtnRB.Y)
	// 		TaoliveOp.MoveClickTitle(closeBtnLT, closeBtnRB)
	// 		robotgo.Sleep(2)
	// 	}
	// }

	for {
		bInIngotCenter, err := TaoliveOp.IsInIngotCenter()
		if err != nil {
			panic(err)
		}

		if !bInIngotCenter {
			fmt.Println("在 iPhone 主界面寻找 点淘 并点击")
			if err := IPhoneOp.GotoTaoLive(); err != nil {
				panic(err)
			}

			fmt.Println("在 点淘 主界面寻找 元宝中心 并点击")
			if err := TaoliveOp.GotoIngotCenter(); err != nil {
				panic(err)
			}
		} else {
			fmt.Println("已在 点淘-元宝中心")
		}

		if TaoliveOp.OCRMoveClickTitle("领取奖励", 10) {
			robotgo.Sleep(3)
			err := ocr.Ocr(nil, nil, nil, nil)
			if err != nil {
				panic(err)
			}
			if TaoliveOp.OCRMoveClickTitle("额外获得68元宝", 0) {
				TaoliveOp.WatchAD("元宝中心", "")
			} else if TaoliveOp.OCRMoveClickTitle("额外获得99元宝", 0) {
				TaoliveOp.WatchAD("元宝中心", "")
			} else if TaoliveOp.OCRMoveClickTitle("看直播60秒得68元宝", 0) {
				TaoliveOp.WatchAD("元宝中心", "")
			} else {
				TaoliveOp.OCRMoveClickTitle("开心收下", 0)
			}
		}

		for TaoliveOp.ExistText("立即领奖") {
			if TaoliveOp.OCRMoveClickTitle("立即领奖", 0) {
				TaoliveOp.WatchAD("元宝中心", "")
				err := ocr.Ocr(nil, nil, nil, nil)
				if err != nil {
					panic(err)
				}
			}
		}

		fmt.Println("在 元宝中心 主界面寻找 睡觉赚元宝 并点击")
		if err := TaoliveOp.GotoSleepToEarn(); err != nil {
			panic(err)
		}

		if err := TaoliveOp.DoSleepToEarn(); err != nil {
			panic(err)
		}

		fmt.Println("在 元宝中心 主界面寻找 下单返元宝 并点击")
		if err := TaoliveOp.GotoOrderToEarn(); err != nil {
			panic(err)
		}

		if err := TaoliveOp.DoOrderToEarn(); err != nil {
			panic(err)
		}

		fmt.Println("在 元宝中心 主界面寻找 今日签到 并点击")
		if err := TaoliveOp.GotoDailySignIn(); err != nil {
			panic(err)
		}

		if err := TaoliveOp.DoDailySignIn(); err != nil {
			panic(err)
		}

		fmt.Println("在 元宝中心 主界面寻找 赚钱卡 并点击")
		if err := TaoliveOp.GotoEarnMoneyCard(); err != nil {
			panic(err)
		}

		if err := TaoliveOp.DoEarnMoneyCard(); err != nil {
			panic(err)
		}

		fmt.Println("在 元宝中心 主界面寻找 走路赚元宝 并点击")
		if err := TaoliveOp.GotoWalkToEarn(); err != nil {
			panic(err)
		}

		if err := TaoliveOp.DoWalkToEarn(); err != nil {
			panic(err)
		}

		fmt.Println("在 元宝中心 主界面寻找 摇一摇赚元宝 并点击")
		if err := TaoliveOp.GotoShakeToEarn(); err != nil {
			panic(err)
		}

		if err := TaoliveOp.DoShakeToEarn(); err != nil {
			panic(err)
		}

		fmt.Println("在 元宝中心 主界面寻找 打工赚元宝 并点击")
		if err := TaoliveOp.GotoWorkToEarn(); err != nil {
			panic(err)
		}

		if err := TaoliveOp.DoWorkToEarn(); err != nil {
			panic(err)
		}
		break
	}
}
