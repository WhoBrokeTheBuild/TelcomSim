package main

import (
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"

	"github.com/WhoBrokeTheBuild/TelcomSim/stbi"

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

	stbi.SetFlipVerticallyOnLoad(stbi.True)

	Loadf("Texture [%v]", filename)
	b, err := LoadAsset(filename)
	if err != nil {
		return err
	}

	image, w, h, _ := stbi.LoadFromMemory(b, stbi.RGBAlpha)
	defer stbi.ImageFree(image)

	gl.GenTextures(1, &t.ID)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(w),
		int32(h),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(image))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return nil
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}
