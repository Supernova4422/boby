package command

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func renderText(text string) (draw.Image, error) {
	content, err := ioutil.ReadFile("Quivira.otf")
	if err != nil {
		return nil, err
	}

	f, err := opentype.Parse(content)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}

	d := &font.Drawer{
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: face,
	}

	// Draw it to get the width needed for the image, then redo it.
	bound, _ := d.BoundString(text)
	d.Dst = image.NewGray(image.Rect(
		bound.Min.X.Floor(),
		bound.Min.Y.Floor(),
		bound.Max.X.Ceil(),
		bound.Max.Y.Ceil()),
	)
	d.DrawString(text)

	return d.Dst, nil
}

func RenderText(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if len(msg) == 0 {
		return
	}

	image, err := renderText(msg[0].(string))
	if err != nil {
		sink(sender, service.Message{Title: "An error has occured"})
	} else {
		file, err := os.Create("hello-go.png")
		if err != nil {
			panic(err)
		}

		defer file.Close()
		if err := png.Encode(file, image); err != nil {
			panic(err)
		}

		sink(
			sender, service.Message{
				Title: "Message received.",
			},
		)
	}
}
