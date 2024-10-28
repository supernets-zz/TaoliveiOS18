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

func Ocr(x, y, w, h *int) ([]interface{}, error) {
	if x == nil {
		x = &AppX
	}

	if y == nil {
		y = &AppY
	}

	if w == nil {
		w = &AppWidth
	}

	if h == nil {
		h = &AppHeight
	}

	fmt.Printf("%s.%03d, Save screenshot\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().UnixMilli()%1000)
	var img image.Image
	var err error
	if *h != AppHeight {
		img, err = robotgo.CaptureImg(*x+10, *y+38, *w-20, *h)
	} else {
		img, err = robotgo.CaptureImg(*x+10, *y+38, *w-20, *h-38-5)
	}
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
	fmt.Printf("%s.%03d, OCR Complete\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().UnixMilli()%1000)
	var result []interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	for _, v := range result {
		txt := v.([]interface{})[1].([]interface{})[0]
		Polygon := v.([]interface{})[0]
		var leftTop, rightBtm robotgo.Point
		leftTop.X = int(Polygon.([]interface{})[0].([]interface{})[0].(float64))
		leftTop.Y = int(Polygon.([]interface{})[0].([]interface{})[1].(float64))
		rightBtm.X = int(Polygon.([]interface{})[2].([]interface{})[0].(float64))
		rightBtm.Y = int(Polygon.([]interface{})[2].([]interface{})[1].(float64))
		fmt.Printf("(%3d, %3d)-(%3d, %3d) %s\n", leftTop.X, leftTop.Y, rightBtm.X, rightBtm.Y, txt)
	}

	return result, nil
}
