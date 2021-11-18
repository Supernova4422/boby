package command

import (
	"bytes"
	"fmt"
	"html"
	"image"
	"image/draw"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// Use inkscape to render text as an image.
func renderText(input string) (io.Reader, error) {
	template := `<svg height="500" width="500">
<text x="0" y="15">%s</text>
</svg>`
	svg := fmt.Sprintf(template, html.EscapeString(input))
	cmd := exec.Command("/usr/bin/inkscape", "--file=-", "--export-png=-", "--export-background=white", "--export-dpi=300")
	buf := new(bytes.Buffer)
	cmd.Stdin = strings.NewReader(svg)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return strings.NewReader(buf.String()), nil
}

// Crop whitespace from the right and bottom of an image.
func cropImage(img image.Image) image.Image {
	breakNow := false
	for x := img.Bounds().Max.X - 1; x > img.Bounds().Min.X; x-- {
		for y := img.Bounds().Max.Y - 1; y > img.Bounds().Min.Y; y-- {
			col := img.At(x, y)
			r, g, b, a := col.RGBA()
			if r != 65535 || g != 65535 || b != 65535 || a != 65535 {
				newBox := image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y, x+5, img.Bounds().Max.Y)
				m := image.NewRGBA(newBox)
				draw.Draw(m, newBox, img, image.Pt(0, 0), draw.Src)
				img = m
				breakNow = true
				break
			}
		}
		if breakNow {
			break
		}
	}

	breakNow = false
	for y := img.Bounds().Max.Y - 1; y > img.Bounds().Min.Y; y-- {
		for x := img.Bounds().Max.X - 1; x > img.Bounds().Min.X; x-- {
			col := img.At(x, y)
			r, g, b, a := col.RGBA()
			if r != 65535 || g != 65535 || b != 65535 || a != 65535 {
				newBox := image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y, img.Bounds().Max.X, y+5)
				m := image.NewRGBA(newBox)
				draw.Draw(m, newBox, img, image.Pt(0, 0), draw.Src)
				img = m
				breakNow = true
				break
			}
		}
		if breakNow {
			break
		}
	}
	return img
}

// RenderText renders a message as an image then replies with the image.
func RenderText(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if len(msg) == 0 {
		return
	}

	png, err := renderText(msg[0].(string))
	if err != nil {
		return
	}

	img, _, err := image.Decode(png)
	if err != nil {
		return
	}

	img = cropImage(img)

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
