package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/WhoBrokeTheBuild/TelcomSim/asset"
	"github.com/WhoBrokeTheBuild/TelcomSim/context"
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Overlay represents a UI layer
type Overlay struct {
	Buffer     *image.RGBA
	Texture    *asset.Texture
	Shader     *asset.Shader
	Mesh       *asset.Mesh
	RenderCtx  context.Render
	Size       mgl32.Vec2
	Components []Component

	needTextureUpdate bool
}

// NewOverlay returns a new Overlay of the given size
func NewOverlay(size mgl32.Vec2) (*Overlay, error) {
	shader, err := asset.NewShaderFromFiles([]string{
		"shaders/ui.vs.glsl",
		"shaders/ui.fs.glsl",
	})
	if err != nil {
		return nil, err
	}

	data := asset.MeshData{
		Vertices: []mgl32.Vec3{
			mgl32.Vec3{size.X(), 0, 0},
			mgl32.Vec3{size.X(), size.Y(), 0},
			mgl32.Vec3{0, size.Y(), 0},
			mgl32.Vec3{size.X(), 0, 0},
			mgl32.Vec3{0, size.Y(), 0},
			mgl32.Vec3{0, 0, 0},
		},
		TexCoords: []mgl32.Vec2{
			mgl32.Vec2{1, 0},
			mgl32.Vec2{1, 1},
			mgl32.Vec2{0, 1},
			mgl32.Vec2{1, 0},
			mgl32.Vec2{0, 1},
			mgl32.Vec2{0, 0},
		},
	}
	mesh, err := asset.NewMeshFromData([]asset.MeshData{data})
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{int(size.X()), int(size.Y())},
	})
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("Failed to create UI image, Unsupported Stride")
	}

	texture, err := asset.NewTextureFromData(rgba.Pix, int(size.X()), int(size.Y()))
	if err != nil {
		return nil, err
	}

	return &Overlay{
		Buffer:  rgba,
		Texture: texture,
		Shader:  shader,
		Mesh:    mesh,
		Size:    size,

		RenderCtx: context.Render{
			Projection: mgl32.Ortho2D(0, size.X(), 0, size.Y()),
		},

		needTextureUpdate: true,
	}, nil
}

// Delete frees all resources owned by the Overlay
func (o *Overlay) Delete() {
	if o.Texture != nil {
		o.Texture.Delete()
		o.Texture = nil
	}
	if o.Shader != nil {
		o.Shader.Delete()
		o.Shader = nil
	}
	if o.Mesh != nil {
		o.Mesh.Delete()
		o.Mesh = nil
	}
}

// Update processes all events
func (o *Overlay) Update(ctx *context.Update) {

}

func (o *Overlay) clear() {
	draw.Draw(o.Buffer, o.Buffer.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.ZP, draw.Src)
}

// Draw renders the current buffer to the screen
func (o *Overlay) Draw() {
	o.clear()
	for _, c := range o.Components {
		c.Draw(o.Buffer)
	}

	o.needTextureUpdate = true
	if o.needTextureUpdate {
		o.updateTexture()
		o.needTextureUpdate = false
	}

	gl.Clear(gl.DEPTH_BUFFER_BIT)

	o.Shader.Bind()

	gl.Uniform1i(o.Shader.GetUniformLocation("uTexture"), 0)
	if o.Buffer != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		o.Texture.Bind()
	} else {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}

	if o.Mesh != nil {
		o.Mesh.Draw(&o.RenderCtx, o.Shader)
	}
}

// AddComponent adds the given Component as a child node of the Overlay
func (o *Overlay) AddComponent(c Component) {
	o.Components = append(o.Components, c)
}

func (o *Overlay) updateTexture() {
	o.Texture.Bind()
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, int32(o.Size.X()), int32(o.Size.Y()), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(o.Buffer.Pix))
}
