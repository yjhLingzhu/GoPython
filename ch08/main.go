package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

// 缩放
func main() {
	f, err := os.Open("D:/Data/zhiyuzhou/portrait_matting/20211220/5FCyWinHXc2afEGi.png")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	gopherImg, _, _ := image.Decode(f) // 打开图片

	// dx := uint(gopherImg.Bounds().Dx() * 3 / 5)
	// dy := uint(gopherImg.Bounds().Dy() * 3 / 5)

	default_width := 295.0  // 尺寸x
	default_height := 413.0 //尺寸y

	img_x := gopherImg.Bounds().Dx() // 图片宽x
	img_y := gopherImg.Bounds().Dy() // 图片高y

	dx := 0.0 // 缩放后的宽x
	dy := 0.0 // 缩放后的高y

	if img_x < int(default_width) || img_y < int(default_height) {
		fmt.Println("上传的图片长或宽过于小！")
		return
	}
	if img_x < int(2*default_width) && img_y > int(2*default_height) {
		fmt.Println("x小于2倍，y大于2倍...")
		dx = float64(gopherImg.Bounds().Dx()) * default_width / float64(img_x)
		dy = float64(gopherImg.Bounds().Dy()) * default_width / float64(img_x)
	} else if img_y < int(2*default_height) && img_x > int(2*default_width) {
		fmt.Println("x大于2倍，y小于2倍...")
		dx = float64(gopherImg.Bounds().Dx()) * default_height / float64(img_y)
		dy = float64(gopherImg.Bounds().Dy()) * default_height / float64(img_y)
	} else if img_x < int(2*default_width) && img_y < int(2*default_height) {
		fmt.Println("直接去人像检测")
		dx = float64(img_x)
		dy = float64(img_y)
	} else {
		fmt.Println("都大于或等于2倍...")
		widthRatio := default_width / float64(img_x)
		heightRatio := default_height / float64(img_y)
		ratio := 0.0
		if widthRatio < heightRatio {
			ratio = widthRatio
		} else {
			ratio = heightRatio
		}
		// fmt.Println("ratio: ", ratio, widthRatio, heightRatio)
		dx = float64(gopherImg.Bounds().Dx()) * ratio
		dy = float64(gopherImg.Bounds().Dy()) * ratio

		// 如果按缩放比缩小后它的长度或者宽度比默认值小的话，那么直接按1/2来缩放  (只针对图片图片长宽是默认值的两倍或以上的生效)
		if dx < default_width || dy < default_height {
			dx = float64((gopherImg.Bounds().Dx() * 1 / 2))
			dy = float64((gopherImg.Bounds().Dy() * 1 / 2))
		}
	}
	// 等比例缩放
	m := resize.Resize(uint(dx), uint(dy), gopherImg, resize.Lanczos3)
	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create("D:/Data/zhiyuzhou/portrait_matting/20211220/scalratio.png") // 创建一张空的图片文件
	defer outFile.Close()
	if err != nil {
		panic(err)
	}
	b := bufio.NewWriter(outFile) // NewWriter返回一个新的Writer，其缓冲区具有默认大小
	err = png.Encode(b, m)        // 将绘制完的图片img重新编码并输入到Writer b中
	if err != nil {
		panic(err)
	}
	err = b.Flush() // 将所有缓冲数据写入底层io.Writer
	if err != nil {
		panic(err)
	}
}
