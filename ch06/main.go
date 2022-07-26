package main

import (
	"bufio"
	"image"
	"image/png"
	_ "image/png"
	"os"
)

// 裁剪
func main() {
	f, err := os.Open("D:/Data/zhiyuzhou/portrait_matting/20211216/gopher2.png")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	gopherImg, _, _ := image.Decode(f) // 打开图片
	img := gopherImg.(*image.RGBA)
	// img := image.NewRGBA(image.Rect(0, 0, 630, 1120))

	subImage := img.SubImage(image.Rect(132, 670, 188, 726)) // 裁剪

	outFile, err := os.Create("D:/Data/zhiyuzhou/portrait_matting/20211216/tailor.png") // 创建一张空的图片文件
	defer outFile.Close()
	if err != nil {
		panic(err)
	}
	b := bufio.NewWriter(outFile) // NewWriter返回一个新的Writer，其缓冲区具有默认大小
	err = png.Encode(b, subImage) // 将绘制完的图片img重新编码并输入到Writer b中
	if err != nil {
		panic(err)
	}
	err = b.Flush() // 将所有缓冲数据写入底层io.Writer
	if err != nil {
		panic(err)
	}
}
