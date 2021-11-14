package command

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/ninetwentyfour/go-wkhtmltoimage"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func getFont() (font.Face, error) {
	content, err := ioutil.ReadFile("noto kufi arabic.ttf")
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

	return face, err
}

func renderText(face font.Face, text string) (draw.Image, error) {
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

// RenderText renders a message as an image then replies with the image.
func RenderText(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if len(msg) == 0 {
		return
	}

	html := `<!DOCTYPE html>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-16">
<html>
    <body>
        <p>%s</p>
    </body>
</html>
`
	options := wkhtmltoimage.ImageOptions{
		BinaryPath: "/usr/local/bin/wkhtmltoimage",
		Input:      "-",
		Format:     "png",
		Html:       fmt.Sprintf(html, msg[0].(string)),
	}
	out, err := wkhtmltoimage.GenerateImage(&options)
	if err != nil {
		log.Fatal(err)
	}

	image, _, err := image.Decode(bytes.NewReader(out))
	if err != nil {
		sink(sender, service.Message{Title: "An error has occured"})
	} else {
		sink(
			sender, service.Message{
				Title: "Rendered text.",
				Image: image,
			},
		)
	}
}
