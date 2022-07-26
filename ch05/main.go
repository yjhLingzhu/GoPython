package main

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"os"
)

// 给png图片换背景
func main() {
	// f, err := os.Open("D:/Data/zhiyuzhou/portrait_matting/20211216/apuhsrnGOIUQQiSe.png")
	// if err != nil {
	// 	panic(err)
	// }
	// img, formatName, err := image.Decode(f)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(formatName)
	// fmt.Println(img.Bounds())
	// fmt.Println(img.ColorModel())

	f, err := os.Open("D:/Data/zhiyuzhou/portrait_matting/20211220/3.png")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	gopherImg, _, _ := image.Decode(f) // 打开图片

	img := image.NewRGBA(image.Rect(0, 0, gopherImg.Bounds().Dx(), gopherImg.Bounds().Dy()))
	for x := 0; x < img.Bounds().Dx(); x++ { // 将背景图涂黑
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, color.RGBA{0, 191, 255, 255}) // 0, 191, 255 天蓝色
		}
	}
	draw.Draw(img, img.Bounds(), gopherImg, image.Pt(0, 0), draw.Over) // 将gopherImg绘制到背景图上
	// subImage := img.SubImage(image.Rect(0, 0, 295, 413))               // 裁剪
	// m := resize.Resize(295, 413, img, resize.Lanczos3)
	// fmt.Println(m)

	outFile, err := os.Create("D:/Data/zhiyuzhou/portrait_matting/20211220/gopher4.png") // 创建一张空的图片文件
	defer outFile.Close()
	if err != nil {
		panic(err)
	}
	b := bufio.NewWriter(outFile) // NewWriter返回一个新的Writer，其缓冲区具有默认大小
	err = png.Encode(b, img)      // 将绘制完的图片img重新编码并输入到Writer b中
	if err != nil {
		panic(err)
	}
	err = b.Flush() // 将所有缓冲数据写入底层io.Writer
	if err != nil {
		panic(err)
	}
}
