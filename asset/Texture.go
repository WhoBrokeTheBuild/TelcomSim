package asset

import (
	"path/filepath"

	"github.com/WhoBrokeTheBuild/TelcomSim/data"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
	"github.com/WhoBrokeTheBuild/TelcomSim/stbi"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Texture represents an OpenGL Texture
type Texture struct {
	ID   uint32
	Size mgl32.Vec2
}

type glTexture struct {
	ID       uint32
	UseCount int
}

var _textures map[string]*glTexture

func init() {
	_textures = map[string]*glTexture{}
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
func NewTextureFromData(data []byte, intFormat uint32, format int32, width, height int) (*Texture, error) {
	t := &Texture{}
	err := t.LoadFromData(data, intFormat, format, width, height)
	if err != nil {
		t.Delete()
		return nil, err
	}
	return t, nil
}

// Delete frees the resources owned by the Texture
func (t *Texture) Delete() {
	if t.ID != InvalidID {
		found := false
		for f, a := range _textures {
			if a.ID == t.ID {
				found = true
				a.UseCount--
				if a.UseCount <= 0 {
					gl.DeleteTextures(1, &a.ID)
					delete(_textures, f)
				}
			}
		}
		if !found {
			gl.DeleteTextures(1, &t.ID)
		}
		t.ID = InvalidID
	}
}

// LoadFromFile loads a Texture from a given file
func (t *Texture) LoadFromFile(filename string) error {
	filename = filepath.Clean(filename)
	t.Delete()

	if a, found := _textures[filename]; found {
		a.UseCount++
		t.ID = a.ID
		log.Loadf("asset.Texture @[%v]", filename)
		return nil
	}

	//stbi.SetFlipVerticallyOnLoad(stbi.True)

	log.Loadf("asset.Texture [%v]", filename)
	b, err := data.Asset(filename)
	if err != nil {
		return err
	}

	image, w, h, ch := stbi.LoadFromMemory(b, stbi.Null)
	defer stbi.ImageFree(image)

	t.Size = mgl32.Vec2{float32(w), float32(h)}

	format := int32(gl.RGB)
	if ch == 4 {
		format = gl.RGBA
	}

	gl.GenTextures(1, &t.ID)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(gl.TEXTURE_2D, 0, format,
		int32(w),
		int32(h),
		0, uint32(format), gl.UNSIGNED_BYTE, gl.Ptr(image))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	if a, found := _textures[filename]; found {
		gl.DeleteTextures(1, &a.ID)
		delete(_textures, filename)
	}

	_textures[filename] = &glTexture{
		ID:       t.ID,
		UseCount: 1,
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return nil
}

// LoadFromData loads a Texture from the given data, width, and height
func (t *Texture) LoadFromData(data []uint8, intFormat uint32, format int32, width, height int) error {
	t.Delete()

	t.Size = mgl32.Vec2{float32(width), float32(height)}

	gl.GenTextures(1, &t.ID)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	if data != nil {
		gl.TexImage2D(gl.TEXTURE_2D, 0, format,
			int32(width),
			int32(height),
			0, uint32(format), gl.UNSIGNED_BYTE, gl.Ptr(data))
	} else {
		gl.TexImage2D(gl.TEXTURE_2D, 0, format,
			int32(width),
			int32(height),
			0, uint32(format), gl.UNSIGNED_BYTE, nil)
	}
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return nil
}

// Bind calls glBindTexture with the Texture's ID
func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}
