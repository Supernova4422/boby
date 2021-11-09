package command

import (
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"golang.org/x/image/math/fixed"
)

// TODO Error handling.
func RenderText(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if len(msg) == 0 {
		return
	}

	content, err := ioutil.ReadFile("Quivira.otf")
	if err != nil {
		panic(err)
	}

	f, err := opentype.Parse(content)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}

	margin := 20 * 64
	margins := fixed.Point26_6{
		X: fixed.Int26_6(margin), // Used to be 20
		Y: fixed.Int26_6(margin), // used to be 30
	}

	d := &font.Drawer{
		Dst:  image.NewRGBA(image.Rect(0, 0, 300, 300)),
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}),
		Face: face,
		Dot:  margins,
	}

	d.DrawString(msg[0].(string))

	d.Dst = image.NewRGBA(image.Rect(0, 0, d.Dot.X.Ceil()+margin, d.Dot.Y.Ceil()+margin))
	d.Dot = margins
	d.DrawString(msg[0].(string))

	file, err := os.Create("hello-go.png")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	if err := png.Encode(file, d.Dst); err != nil {
		panic(err)
	}

	sink(
		sender, service.Message{
			Title: "Message received.",
		},
	)
}
