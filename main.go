package main

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	fontKai *truetype.Font // 字体
	fontTtf *truetype.Font // 字体
)

func main() {
	// 根据路径打开模板文件
	templateFile, err := os.Open("header.png")
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
	fontKai, err = loadFont("MSYHBD.TTC")
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	fontTtf, err = loadFont("MSYHBD.TTC")
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	// 向图片中写入文字
	writeWord2Pic(newTemplateImage)

	// ====================向模板中粘贴图片 begin========================
	// 		1、根据地址获取图片内容
	imageData, err := getDataByUrl("http://qiniu.yueda.vip/123.png")
	if err != nil {
		fmt.Println("根据地址获取图片失败,err:", err.Error())
		return
	}
	// 图片到边框距离
	pic2FramePadding := 20
	// 获取全景图原始的尺寸
	dx := imageData.Bounds().Dx()
	dy := imageData.Bounds().Dy()
	// 		2、重新调整要粘贴图片尺寸
	if dx > dy { // 判断是横图还是竖图
		imageData = resize.Resize(uint(387-pic2FramePadding), uint(183-pic2FramePadding), imageData, resize.Lanczos3)
	} else {
		imageData = resize.Resize(uint(387/2-pic2FramePadding), uint(183-pic2FramePadding), imageData, resize.Lanczos3)
	}

	// 新建一个透明图层
	transparentImg := image.NewRGBA(image.Rect(0, 0, imageData.Bounds().Dx()+pic2FramePadding, imageData.Bounds().Dy()+pic2FramePadding))
	// 将缩略图放到透明图层上
	draw.Draw(transparentImg,
		image.Rect(pic2FramePadding/2, pic2FramePadding/2, transparentImg.Bounds().Dx(), transparentImg.Bounds().Dy()),
		imageData,
		image.Point{},
		draw.Over)

	// 图片周围画线
	lineToPic(transparentImg)

	// 	粘贴缩略图
	draw.Draw(newTemplateImage,
		transparentImg.Bounds().Add(image.Pt(228, 558)),
		transparentImg,
		transparentImg.Bounds().Min,
		draw.Over)

	// // ====================向模板中粘贴图片 结束========================

	// // 保存图片  ---> 在此我们统一将文件保存到：C:\Users\yida\GolandProjects\GoProjectDemo\A-go-study\dst.png
	saveFile(newTemplateImage)
}
func lineToPic(transparentImg *image.RGBA) {
	gc := draw2dimg.NewGraphicContext(transparentImg)
	gc.SetStrokeColor(color.RGBA{ // 线框颜色
		R: uint8(36),
		G: uint8(106),
		B: uint8(96),
		A: 0xff})
	gc.SetFillColor(color.RGBA{})
	gc.SetLineWidth(5) // 线框宽度
	gc.BeginPath()
	gc.MoveTo(0, 0)
	gc.LineTo(float64(transparentImg.Bounds().Dx()), 0)
	gc.LineTo(float64(transparentImg.Bounds().Dx()), float64(transparentImg.Bounds().Dy()))
	gc.LineTo(0, float64(transparentImg.Bounds().Dy()))
	gc.LineTo(0, 0)
	gc.Close()
	gc.FillStroke()
}

func writeWord2Pic(newTemplateImage *image.RGBA) {
	// 在写入之前有一些准备工作
	content := freetype.NewContext()
	content.SetClip(newTemplateImage.Bounds())
	content.SetDst(newTemplateImage)
	content.SetSrc(image.Black) // 设置字体颜色
	content.SetDPI(72)          // 设置字体分辨率

	content.SetFontSize(40)  // 设置字体大小
	content.SetFont(fontKai) // 设置字体样式，就是我们上面加载的字体

	// 	正式写入文字
	// 参数1：要写入的文字
	// 参数2：文字坐标
	content.DrawString("yida同志:", freetype.Pt(150, 375))
	content.DrawString("您在2022年度中表现突出，忠诚奉献、认真负责，", freetype.Pt(200, 450))
	content.DrawString("被评为", freetype.Pt(200, 520))

	content.DrawString("特发此证，以资鼓励。", freetype.Pt(630, 520))
	// 设置字体大小
	content.SetFontSize(48)
	// 设置字体颜色
	content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255}))
	content.DrawString("“最佳奉献奖”", freetype.Pt(300, 520))

	content.SetFont(fontTtf) // 设置字体样式
	// 设置字体大小
	content.SetFontSize(32)
	// 设置字体颜色
	content.SetSrc(image.Black)
	content.DrawString("东风战略导弹部队", freetype.Pt(898, 660))
	content.DrawString("二零二二年五月", freetype.Pt(898, 726))
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
	dstFile, err := os.Create("dst.png")
	if err != nil {
		fmt.Println(err)
	}
	defer dstFile.Close()
	png.Encode(dstFile, pic)
}

// 根据地址获取图片内容
func getDataByUrl(url string) (img image.Image, err error) {
	res, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("[%s]通过url获取数据失败,err:%s", url, err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	// 读取获取的[]byte数据
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取数据失败,err:%s", err.Error())
		return
	}

	if !strings.HasSuffix(url, ".jpg") &&
		!strings.HasSuffix(url, ".jpeg") &&
		!strings.HasSuffix(url, ".png") {
		err = fmt.Errorf("[%s]不支持的图片类型,暂只支持.jpg、.png文件类型", url)
		return
	}

	// []byte 转 io.Reader
	reader := bytes.NewReader(data)
	if strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".jpeg") {
		// 此处jgeg.decode 有坑，明明是.jpg的图片但 会报 invalid JPEG format: missing SOI marker 错误
		// 所以当报错时我们再用 png.decode 试试
		img, err = jpeg.Decode(reader)
		if err != nil {
			fmt.Printf("jpeg.Decode err:%s", err.Error())
			reader2 := bytes.NewReader(data)
			img, err = png.Decode(reader2)
			if err != nil {
				err = fmt.Errorf("===>png.Decode err:%s", err.Error())
				return
			}
		}
	}

	if strings.HasSuffix(url, ".png") {
		img, err = png.Decode(reader)
		if err != nil {
			err = fmt.Errorf("png.Decode err:%s", err.Error())
			return
		}
	}

	return
}
