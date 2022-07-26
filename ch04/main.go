package main

import (
	"fmt"
	"time"
)

// "gocv.io/x/gocv"

// 目标检测 gocv
func main() {
	// fmt.Println(123)
	// dst := "D:/Data/zhiyuzhou/portrait_matting/20211216/gopher2.png"
	// img := gocv.IMRead(dst, -1)
	// defer img.Close()

	// classifier := gocv.NewCascadeClassifier()
	// defer classifier.Close()

	// if !classifier.Load("D:/workspace/owner/go/studyPrograms/src/gocv.io/x/gocv/data/haarcascade_frontalface_default.xml") {
	// 	fmt.Printf("Error reading cascade file: %v\n", "xmlFile")
	// 	return
	// }

	// rects := classifier.DetectMultiScale(img)
	// fmt.Println(rects)

	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)

	remainSecond := time.Duration(todayLastTime.Unix()-time.Now().Local().Unix()) * time.Second

	fmt.Println(remainSecond)
	fmt.Println(remainSecond.Seconds())
}
