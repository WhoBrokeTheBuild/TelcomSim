package asset

import (
	"path/filepath"

	"github.com/WhoBrokeTheBuild/TelcomSim/data"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
	"github.com/WhoBrokeTheBuild/TelcomSim/stbi"

	gl "github.com/go-gl/gl/v4.1-core/gl"
)

// Texture represents an OpenGL Texture
type Texture struct {
	ID uint32
}

// NewTextureFromFile returns a new Texture from the given file
func NewTextureFromFile(filename string) (*Texture, error) {
	t := &Texture{}
	err := t.LoadFromFile(filename)
	if err != nil {
		t.Delete()
		return nil, err
	}
	return t, nil
}

// NewTextureFromData returns a new Texture from the given data, width, and height
func NewTextureFromData(data []byte, width, height int) (*Texture, error) {
	t := &Texture{}
	err := t.LoadFromData(data, width, height)
	if err != nil {
		t.Delete()
		return nil, err
	}
	return t, nil
}

// Delete frees the resources owned by the Texture
func (t *Texture) Delete() {
	if t.ID != InvalidID {
		gl.DeleteTextures(1, &t.ID)
		t.ID = InvalidID
	}
}

// LoadFromFile loads a Texture from a given file
func (t *Texture) LoadFromFile(filename string) error {
	filename = filepath.Clean(filename)
	t.Delete()

	stbi.SetFlipVerticallyOnLoad(stbi.True)

	log.Loadf("asset.Texture [%v]", filename)
	b, err := data.Asset(filename)
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

// LoadFromData loads a Texture from the given data, width, and height
func (t *Texture) LoadFromData(data []byte, width, height int) error {
	t.Delete()

	gl.GenTextures(1, &t.ID)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(width),
		int32(height),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return nil
}

// Bind calls glBindTexture with the Texture's ID
func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}
