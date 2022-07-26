package main

// 人脸检测，基于pico方法

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"

	pigo "github.com/esimov/pigo/core"
	"gocv.io/x/gocv"
)

func main() {

	dst := "D:/Data/zhiyuzhou/portrait_matting/20211216/gopher2.png"
	img := gocv.IMRead(dst, -1)
	defer img.Close()

	window := gocv.NewWindow("pigo")
	defer window.Close()

	window.ResizeWindow(640, 480)

	green := color.RGBA{0, 255, 0, 0}

	cascadeFile, err := ioutil.ReadFile("D:/workspace/owner/go/studyPrograms/src/gocv.io/x/gocv/data/haarcascade_frontalface_default.xml")

	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}

	for {
		goImg, err := img.ToImage()

		if nil != err {
			fmt.Println("ToImage err ")
			return
		}

		pixels := pigo.RgbToGrayscale(goImg)
		cols, rows := goImg.Bounds().Max.X, goImg.Bounds().Max.Y

		cParams := pigo.CascadeParams{
			MinSize:     20,
			MaxSize:     1200,
			ShiftFactor: 0.1,
			ScaleFactor: 1.1,

			ImageParams: pigo.ImageParams{
				Pixels: pixels,
				Rows:   rows,
				Cols:   cols,
				Dim:    cols,
			},
		}

		pPigo := pigo.NewPigo()

		classifier, err := pPigo.Unpack(cascadeFile)
		if err != nil {
			log.Fatalf("Error reading the cascade file: %s", err)
		}

		angle := 0.0
		iouThreshold := 0.3

		dets := classifier.RunCascade(cParams, angle)

		dets = classifier.ClusterDetections(dets, iouThreshold)

		for _, face := range dets {

			if face.Q > 5 {
				x := face.Col - face.Scale/2
				y := face.Row - face.Scale/2
				r := image.Rect(x, y, x+face.Scale, y+face.Scale)
				gocv.Rectangle(&img, r, green, 3)

			} else {
				continue
			}

		}

		window.IMShow(img)

		if 27 == window.WaitKey(1) {
			break
		}

	}

}
