package ui

import (

	// JPEG support
	_ "image/jpeg"

	// PNG support
	_ "image/png"

	"github.com/WhoBrokeTheBuild/TelcomSim/asset"
	"github.com/WhoBrokeTheBuild/TelcomSim/context"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Image is a Component that draws an image to the screen
type Image struct {
	BaseComponent
	Bounds  mgl32.Vec4
	Texture *asset.Texture
	Mesh    *asset.Mesh
}

// NewImageFromFile returns a new Image from the given file
func NewImageFromFile(filename string) *Image {
	c := &Image{}
	err := c.LoadFromFile(filename)
	if err != nil {
		c.Delete()
		log.Errorf("%v", err)
		return nil
	}
	return c
}

// NewImageFromData returns a new Image from the given data, width, and height
func NewImageFromData(data []uint8, intFormat uint32, format int32, width, height int) *Image {
	c := &Image{}
	err := c.LoadFromData(data, intFormat, format, width, height)
	if err != nil {
		c.Delete()
		log.Errorf("%v", err)
		return nil
	}
	return c
}

// Delete frees all resources owned by the Image
func (c *Image) Delete() {
	if c.Texture != nil {
		c.Texture.Delete()
		c.Texture = nil
	}
}

// LoadFromFile loads an Image from the given file
func (c *Image) LoadFromFile(filename string) error {
	var err error
	c.Delete()

	c.Texture, err = asset.NewTextureFromFile(filename)
	if err != nil {
		c.Delete()
		return err
	}

	c.SetSize(c.Texture.Size)

	return nil
}

// LoadFromData loads an Image from the given data, width, and height
func (c *Image) LoadFromData(data []uint8, intFormat uint32, format int32, width, height int) error {
	var err error
	c.Delete()

	c.Texture, err = asset.NewTextureFromData(data, intFormat, format, width, height)
	if err != nil {
		c.Delete()
		return err
	}

	c.SetSize(c.Texture.Size)

	return nil
}

// SetPosition sets the Image's position
func (c *Image) SetPosition(pos mgl32.Vec2) {
	c.BaseComponent.SetPosition(pos)
	c.updateMesh()
}

// SetSize sets the Image's size
func (c *Image) SetSize(size mgl32.Vec2) {
	c.BaseComponent.SetSize(size)
	c.updateMesh()
}

func (c *Image) updateMesh() {
	var err error
	pos := c.GetPosition()
	size := c.GetSize()

	x := pos.X()
	y := pos.Y()
	w := size.X()
	h := size.Y()

	if c.Mesh == nil {
		c.Mesh, err = new2DMesh(
			mgl32.Vec4{x, y, x + w, y + h},
			mgl32.Vec4{0, 0, 1, 1})
		if err != nil {
			c.Delete()
			log.Errorf("%v", err)
		}
	} else {
		err = update2DMesh(c.Mesh,
			mgl32.Vec4{x, y, x + w, y + h},
			mgl32.Vec4{0, 0, 1, 1})
		if err != nil {
			c.Delete()
			log.Errorf("%v", err)
		}
	}
}

// Draw renders the Image to the buffer
func (c *Image) Draw(ctx *context.Render) {
	gl.Uniform1i(ctx.Shader.GetUniformLocation("uTexture"), 0)

	gl.ActiveTexture(gl.TEXTURE0)
	if c.Texture != nil {
		c.Texture.Bind()
	}
	if c.Mesh != nil {
		c.Mesh.Draw(ctx)
	}
}
