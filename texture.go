package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"

	gl "github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	ID uint32
}

func NewTexture(filename string) (*Texture, error) {
	t := &Texture{}
	err := t.Load(filename)
	if err != nil {
		t.Delete()
		return nil, err
	}
	return t, nil
}

func (t *Texture) Delete() {
	if t.ID != InvalidID {
		gl.DeleteTextures(1, &t.ID)
		t.ID = InvalidID
	}
}

func (t *Texture) Load(filename string) error {
	filename = filepath.Clean(filename)
	t.Delete()

	Loadf("Texture [%v]", filename)
	b, err := LoadAsset(filename)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return fmt.Errorf("Unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	gl.GenTextures(1, &t.ID)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return nil
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}
