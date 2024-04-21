package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/xyproto/palgen"
)

type ResizeType int

const (
	Vertical ResizeType = iota
	Horizontal
)

type PNGImg struct {
	FileName string
	Img      []byte
}

// args[0]: fileCount, args[1...]: files
func Convert(this js.Value, args []js.Value) interface{} {
	offset := 1

	var tmp = args[0].String()
	fileCount, err := strconv.Atoi(tmp)
	if err != nil {
		printAlert("添付されたファイルカウントの取得に失敗しました")
		return nil
	}

	var delays []int
	var images []*image.Paletted

	for i := offset; i < offset+fileCount; i++ {
		var num = args[i].Get("delay").String()
		delay, err := strconv.Atoi(num)
		if err != nil {
			printAlert("時間の取得に失敗しました")
			return nil
		}

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
		pal, err := palgen.Generate(img, 256)
		palettedImg := image.NewPaletted(img.Bounds(), pal)
		draw.FloydSteinberg.Draw(palettedImg, palettedImg.Rect, img, img.Bounds().Min)
		images = append(images, palettedImg)

		delays = append(delays, delay)
	}

	var gifData bytes.Buffer
	gif.EncodeAll(&gifData, &gif.GIF{Image: images, Delay: delays})

	attachData(gifData.Bytes(), "output", ".gif")
	return nil
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
