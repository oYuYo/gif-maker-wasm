package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
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
	delay, err := strconv.Atoi(tmp)
	if err != nil {
		fmt.Println(err)
		printAlert("時間の取得に失敗しました")
		return nil
	}

	tmp = args[1].String()
	fileCount, err := strconv.Atoi(tmp)
	if err != nil {
		printAlert("添付されたファイルカウントの取得に失敗しました")
		return nil
	}

	var delays []int
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
		palettedImg := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(palettedImg, palettedImg.Bounds(), img, image.Point{}, draw.Over)
		images = append(images, palettedImg)

		delays = append(delays, delay)
	}

	var gifData bytes.Buffer
	gif.EncodeAll(&gifData, &gif.GIF{Image: images, Delay: delays})

	_, err = createGIF(&JPGImgs)
	if err != nil {
		printAlert("zipファイル作成中にエラーが発生しました")
		return nil
	}

	attachData(gifData.Bytes(), "output", ".gif")
	return nil
}

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
