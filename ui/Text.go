package ui

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/WhoBrokeTheBuild/TelcomSim/data"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var _fonts map[string]*truetype.Font

func init() {
	_fonts = map[string]*truetype.Font{}
}

// Text is a Component that draws text to the screen
type Text struct {
	BaseComponent

	Text  string
	Size  float64
	Color color.Color
	Font  *truetype.Font
}

// NewText returns a new Text from a given string, font, font size, and color
func NewText(text string, font string, size float64, color color.Color) *Text {
	var f *truetype.Font

	if tmp, found := _fonts[font]; found {
		log.Loadf("ui.Font [%v]+", font)
		f = tmp
	} else {
		log.Loadf("ui.Font [%v]", font)
		b, err := data.Asset(font)
		if err != nil {
			log.Errorf("%v", err)
			return nil
		}

		f, err = freetype.ParseFont(b)
		if err != nil {
			log.Errorf("%v", err)
			return nil
		}
		_fonts[font] = f
	}

	return &Text{
		Text:  text,
		Size:  size,
		Color: color,
		Font:  f,
	}
}

func (c *Text) SetText(text string) {
	c.Text = text
}

// Draw renders the Text to the buffer
func (c *Text) Draw(buffer draw.Image) {
	pos := c.GetPosition()
	pt := fixed.Point26_6{
		X: fixed.Int26_6(pos.X * 64),
		Y: fixed.Int26_6((pos.Y + 18) * 64),
	}
	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(c.Font)
	ctx.SetFontSize(c.Size)
	ctx.SetClip(buffer.Bounds())
	ctx.SetDst(buffer)
	ctx.SetSrc(image.NewUniform(c.Color))
	ctx.SetHinting(font.HintingFull)
	ctx.DrawString(c.Text, pt)
}
