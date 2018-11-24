package ui

import (
	"bytes"
	"image"
	"image/draw"

	// JPEG support
	_ "image/jpeg"

	// PNG support
	_ "image/png"

	"github.com/WhoBrokeTheBuild/TelcomSim/data"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
)

var _images map[string]*image.RGBA

func init() {
	_images = map[string]*image.RGBA{}
}

// Image is a Component that draws an image to the screen
type Image struct {
	BaseComponent

	RGBA *image.RGBA
}

// NewImageFromFile returns a new Image from the given file
func NewImageFromFile(filename string, bounds image.Rectangle) *Image {
	c := &Image{}
	err := c.LoadFromFile(filename, bounds)
	if err != nil {
		c.Delete()
		log.Errorf("%v", err)
		return nil
	}
	return c
}

// NewImageFromData returns a new Image from the given data, width, and height
func NewImageFromData(data []uint8, width, height int) *Image {
	c := &Image{}
	err := c.LoadFromData(data, width, height)
	if err != nil {
		c.Delete()
		log.Errorf("%v", err)
		return nil
	}
	return c
}

// Delete frees all resources owned by the Image
func (c *Image) Delete() {
}

// LoadFromFile loads an Image from the given file
func (c *Image) LoadFromFile(filename string, bounds image.Rectangle) error {
	var rgba *image.RGBA

	if tmp, found := _images[filename]; found {
		log.Loadf("ui.Image [%v]+", filename)
		rgba = tmp
	} else {
		log.Loadf("ui.Image [%v]", filename)
		b, err := data.Asset(filename)
		if err != nil {
			return err
		}

		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return err
		}

		rgba = image.NewRGBA(img.Bounds())
		draw.Draw(rgba, img.Bounds(), img, image.ZP, draw.Src)
		_images[filename] = rgba
	}

	if bounds.Empty() {
		c.RGBA = rgba
	} else {
		log.Infof("Clipping to %v", bounds)
		c.RGBA = rgba.SubImage(bounds).(*image.RGBA)
	}

	c.SetSize(c.RGBA.Bounds().Size())

	return nil
}

// LoadFromData loads an Image from the given data, width, and height
func (c *Image) LoadFromData(data []uint8, width, height int) error {
	c.RGBA = &image.RGBA{
		Pix:    data,
		Stride: 4 * width,
		Rect:   image.Rectangle{image.ZP, image.Pt(width, height)},
	}
	return nil
}

func (c *Image) bounds() image.Rectangle {
	return image.Rectangle{
		c.GetPosition().Sub(c.RGBA.Bounds().Min),
		c.GetPosition().Add(c.RGBA.Bounds().Max),
	}
}

// Draw renders the Image to the buffer
func (c *Image) Draw(buffer draw.Image) {
	draw.Draw(buffer, c.bounds(), c.RGBA, image.ZP, draw.Src)
}
