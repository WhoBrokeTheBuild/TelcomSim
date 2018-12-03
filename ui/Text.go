package ui

import (
	"image"
	"image/color"

	"github.com/WhoBrokeTheBuild/TelcomSim/asset"
	"github.com/WhoBrokeTheBuild/TelcomSim/data"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var _fonts map[string]*truetype.Font

func init() {
	_fonts = map[string]*truetype.Font{}
}

// Text is a Component that draws text to the screen
type Text struct {
	Image

	Text  string
	Size  float64
	Color color.Color
	Font  *truetype.Font
	Face  font.Face
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

	c := &Text{
		Text:  text,
		Size:  size,
		Color: color,
		Font:  f,
	}
	c.updateTexture()

	return c
}

// SetText sets the text to be rendered
func (c *Text) SetText(text string) {
	c.Text = text
	c.updateTexture()
}

func (c *Text) updateTexture() {
	var err error

	if c.Texture != nil {
		c.Texture.Delete()
	}

	d := &font.Drawer{
		Dst: nil,
		Src: image.NewUniform(c.Color),
		Dot: freetype.Pt(0, int(c.Size)),
		Face: truetype.NewFace(c.Font, &truetype.Options{
			Size:    c.Size,
			DPI:     72,
			Hinting: font.HintingFull,
		}),
	}

	m := d.MeasureString(c.Text)
	buffer := image.NewRGBA(image.Rect(0, 0, m.Ceil(), int(c.Size)))

	d.Dst = buffer
	d.DrawString(c.Text)

	s := buffer.Rect.Size()
	c.Texture, err = asset.NewTextureFromData(buffer.Pix, gl.RGBA, gl.RGBA, s.X, s.Y)
	if err != nil {
		log.Errorf("%v", err)
	}

	c.SetSize(mgl32.Vec2{float32(s.X), float32(s.Y)})
}
