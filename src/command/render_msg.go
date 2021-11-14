package command

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/ninetwentyfour/go-wkhtmltoimage"
	"golang.org/x/image/font"
)

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
	// TODO ENSURE wkhtmltoimage exists!
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

	break_now := false
	img, _, err := image.Decode(bytes.NewReader(out))
	for x := img.Bounds().Max.X - 1; x > img.Bounds().Min.X; x-- {
		for y := img.Bounds().Max.Y - 1; y > img.Bounds().Min.Y; y-- {
			col := img.At(x, y)
			r, g, b, a := col.RGBA()
			if r != 65535 || g != 65535 || b != 65535 || a != 65535 {
				new_box := image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y, x+5, img.Bounds().Max.Y)
				m := image.NewRGBA(new_box)
				draw.Draw(m, new_box, img, image.Pt(0, 0), draw.Src)
				img = m
				break_now = true
				break
			}
		}
		if break_now {
			break
		}
	}
	if err != nil {
		sink(sender, service.Message{Title: "An error has occured"})
	} else {
		sink(
			sender, service.Message{
				Title: "Rendered text.",
				Image: img,
			},
		)
	}
}
