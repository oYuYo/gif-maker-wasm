package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"
	"syscall/js"
)

type ResizeType int

const (
	Vertical ResizeType = iota
	Horizontal
)

type JPGImg struct {
	FileName string
	Img      []byte
}

type PNGImg struct {
	FileName string
	Img      []byte
}

// args[0]: delay, args[1]: fileCount, args[2...]: files
func Convert(this js.Value, args []js.Value) interface{} {
	offset := 2
	var tmp string = args[0].String()
	_, err := strconv.Atoi(tmp)
	if err != nil {
		fmt.Println(err)
		printAlert("秒数の取得に失敗しました")
		return nil
	}

	tmp = args[1].String()
	fileCount, err := strconv.Atoi(tmp)
	if err != nil {
		printAlert("添付されたファイルカウントの取得に失敗しました")
		return nil
	}

	var JPGImgs []JPGImg
	var images []*image.Paletted

	for i := offset; i < offset+fileCount; i++ {
		base64Decode, err := base64.StdEncoding.DecodeString(args[i].Get("base64").String())
		if err != nil {
			printAlert("添付された画像データの取得に失敗しました")
			return nil
		}

		img, err := png.Decode(strings.NewReader(string(base64Decode)))
		if err != nil {
			printAlert("PNGデータのデコードに失敗しました")
			return nil
		}
		palettedImg := image.NewPaletted(img.Bounds(), palette.WebSafe)
		//palettedImg := image.NewPaletted(img.Bounds(), nil)
		//draw.Draw(palettedImg, palettedImg.Bounds(), img, image.Point{}, draw.Over)
		images = append(images, palettedImg)

		var b bytes.Buffer
		if err := jpeg.Encode(bufio.NewWriter(&b), img, &jpeg.Options{Quality: 90}); err != nil {
			printAlert("JPGデータへのエンコードに失敗しました")
			return nil
		}
	}

	var gifData bytes.Buffer
	gif.EncodeAll(&gifData, &gif.GIF{Image: images, Delay: []int{100}})

	_, err = createGIF(&JPGImgs)
	if err != nil {
		printAlert("zipファイル作成中にエラーが発生しました")
		return nil
	}

	attachData(gifData.Bytes(), "output", ".gif")
	return nil
}

/*
func fillTransparentWhite(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	white := color.RGBA{255, 255, 255, 255}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, alpha := img.At(x, y).RGBA()
			if alpha == 0 {
				newImg.Set(x, y, white)
			} else {
				newImg.Set(x, y, img.At(x, y))
			}
		}
	}
	return newImg
}

func resizeImage(resizeType ResizeType, specifiedfileSize int, img image.Image) image.Image {
	origBounds := img.Bounds()
	origWidth := origBounds.Dx()
	origHeight := origBounds.Dy()
	var ratio float64

	if resizeType == Horizontal {
		ratio = float64(specifiedfileSize) / float64(origWidth)
	} else {
		ratio = float64(specifiedfileSize) / float64(origHeight)
	}
	newWidth := int(float64(origWidth) * ratio)
	newHeight := int(float64(origHeight) * ratio)

	resizeImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			origX := int(float64(x) / ratio)
			origY := int(float64(y) / ratio)
			resizeImg.Set(x, y, img.At(origX, origY))
		}
	}

	return resizeImg
}
*/

func createGIF(data *[]JPGImg) ([]byte, error) {
	var gifData bytes.Buffer

	return gifData.Bytes(), nil
}

func attachData(data []byte, fileName string, ext string) {
	document := js.Global().Get("document")
	el := document.Call("getElementById", "output-file")
	encode := base64.StdEncoding.EncodeToString(data)
	dataUri := fmt.Sprintf("data:%s;base64,%s", "image/gif", encode)
	el.Set("href", dataUri)
	el.Set("download", fileName+ext)
	el.Call("click")
}

func printAlert(msg string) {
	document := js.Global().Get("document")
	el := document.Call("getElementById", "err-msg-spn")
	el.Set("innerText", msg)
}

func main() {
	ch := make(chan struct{}, 0)
	js.Global().Set("Convert", js.FuncOf(Convert))
	<-ch
}
