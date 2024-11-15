package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

var AppX, AppY, AppWidth, AppHeight int
var OCRResult []interface{}
var OCRTick int64

func Ocr(x, y, w, h *int) error {
	if x == nil {
		x = &AppX
	} else {
		fmt.Printf("x = %d, ", *x)
	}

	if y == nil {
		y = &AppY
	} else {
		fmt.Printf("y = %d, ", *y)
	}

	if w == nil {
		w = &AppWidth
	} else {
		fmt.Printf("w = %d, ", *w)
	}

	if h == nil {
		h = &AppHeight
	} else {
		fmt.Printf("h = %d\n", *h)
	}

	fmt.Printf("%s.%03d, Save screenshot\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().UnixMilli()%1000)
	OCRTick = time.Now().Unix()
	var img image.Image
	var err error
	img, err = robotgo.CaptureImg(*x, *y, *w, *h)
	if err != nil {
		log.Fatalf("Failed to CaptureImg: %v", err)
	}

	err = robotgo.SavePng(img, "./t.png")
	if err != nil {
		log.Fatalf("Failed to SavePng: %v", err)
	}

	fmt.Printf("%s.%03d, OCRing...\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().UnixMilli()%1000)
	imagePath := "./t.png"
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", imagePath)
	if err != nil {
		log.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalf("Failed to copy file: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "http://localhost:9527/ocr", body)
	// req, err := http.NewRequest("POST", "http://192.168.1.78:9527/ocr", body)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// fmt.Printf("%s\n", respBody)
	k := bytes.IndexByte(respBody, '<')
	if k != -1 {
		respBody = append(respBody[0:k], respBody[k+1:]...)
	}
	k = bytes.IndexByte(respBody, '>')
	if k != -1 {
		respBody = append(respBody[0:k], respBody[k+1:]...)
	}

	fmt.Printf("%s.%03d, OCR Complete\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().UnixMilli()%1000)
	err = json.Unmarshal(respBody, &OCRResult)
	if err != nil {
		fmt.Printf("%s\n", respBody)
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	for _, v := range OCRResult {
		txt := v.([]interface{})[1].([]interface{})[0]
		Polygon := v.([]interface{})[0]
		var leftTop, rightBtm robotgo.Point
		leftTop.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
		leftTop.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
		rightBtm.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
		rightBtm.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
		fmt.Printf("(%3d, %3d)-(%3d, %3d) %s\n", leftTop.X, leftTop.Y, rightBtm.X, rightBtm.Y, txt)
	}

	return nil
}
