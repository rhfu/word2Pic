package main

import (
	// "bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	// "github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/draw"
	// "image/jpeg"
	"image/png"
	// "io"
	"io/ioutil"
	"log"
	// "net/http"
	"os"
	// "strings"
)

var (
	yh *truetype.Font // 字体
)

func main() {
	// 根据路径打开模板文件
	templateFile, err := os.Open("header-2.png")
	if err != nil {
		panic(err)
	}
	defer templateFile.Close()
	// 解码
	templateFileImage, err := png.Decode(templateFile)
	if err != nil {
		panic(err)
	}
	// 新建一张和模板文件一样大小的画布
	newTemplateImage := image.NewRGBA(templateFileImage.Bounds())
	// 将模板图片画到新建的画布上
	draw.Draw(newTemplateImage, templateFileImage.Bounds(), templateFileImage, templateFileImage.Bounds().Min, draw.Over)
	// 加载字体文件  这里我们加载两种字体文件
	yh, err = loadFont("MSYHBD.TTC")
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	// 向图片中写入文字
	writeWord2Pic(newTemplateImage)
	saveFile(newTemplateImage)
}

func writeWord2Pic(newTemplateImage *image.RGBA) {
	// 在写入之前有一些准备工作
	content := freetype.NewContext()
	content.SetClip(newTemplateImage.Bounds())
	content.SetDst(newTemplateImage)
	// content.SetSrc(image.Black) // 设置字体颜色
	content.SetSrc(image.NewUniform(color.RGBA{R: 160, G: 118, B: 93, A: 255}))
	content.SetDPI(50)       // 设置字体分辨率
	content.SetFontSize(100) // 设置字体大小
	content.SetFont(yh)      // 设置字体样式，就是我们上面加载的字体
	content.DrawString("第一单元厚重的和谐文化", freetype.Pt(100, 175))

	content.SetDPI(50)
	content.SetFontSize(100)    // 设置字体大小
	content.SetSrc(image.White) // 设置字体颜色
	content.DrawString("第一单元厚重的和谐文化", freetype.Pt(93, 165))
}

// 根据路径加载字体文件
// path 字体的路径
func loadFont(path string) (font *truetype.Font, err error) {
	var fontBytes []byte
	fontBytes, err = ioutil.ReadFile(path) // 读取字体文件
	if err != nil {
		err = fmt.Errorf("加载字体文件出错:%s", err.Error())
		return
	}
	font, err = freetype.ParseFont(fontBytes) // 解析字体文件
	if err != nil {
		err = fmt.Errorf("解析字体文件出错,%s", err.Error())
		return
	}
	return
}

func saveFile(pic *image.RGBA) {
	dstFile, err := os.Create("dst-2.png")
	if err != nil {
		fmt.Println(err)
	}
	defer dstFile.Close()
	png.Encode(dstFile, pic)
}
