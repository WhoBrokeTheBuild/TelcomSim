package ui

import (
	"fmt"

	"github.com/WhoBrokeTheBuild/TelcomSim/asset"
	"github.com/WhoBrokeTheBuild/TelcomSim/context"
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Overlay represents a UI layer
type Overlay struct {
	Buffer     *asset.Texture
	Shader     *asset.Shader
	Mesh       *asset.Mesh
	RenderCtx  context.Render
	Size       mgl32.Vec2
	Components []Component

	frameID           uint32
	depthID           uint32
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

	mesh, err := new2DMesh(mgl32.Vec4{0, 0, size.X(), size.Y()}, mgl32.Vec4{0, 1, 1, 0})
	if err != nil {
		return nil, err
	}

	var (
		frameID uint32
		depthID uint32
	)

	buffer, err := asset.NewTextureFromData(nil, gl.RGBA8, gl.RGBA, int(size.X()), int(size.Y()))
	if err != nil {
		return nil, err
	}

	gl.GenRenderbuffers(1, &depthID)
	gl.BindRenderbuffer(gl.RENDERBUFFER, depthID)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, int32(size.X()), int32(size.Y()))
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	gl.GenFramebuffers(1, &frameID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameID)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, buffer.ID, 0)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, depthID)

	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		return nil, fmt.Errorf("Failed to create Framebuffer")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return &Overlay{
		Buffer: buffer,
		Shader: shader,
		Mesh:   mesh,
		Size:   size,

		RenderCtx: context.Render{
			Projection: mgl32.Ortho2D(0, size.X(), 0, size.Y()),
			Shader:     shader,
		},

		frameID: frameID,
		depthID: depthID,

		needTextureUpdate: true,
	}, nil
}

// Delete frees all resources owned by the Overlay
func (o *Overlay) Delete() {
	if o.Buffer != nil {
		o.Buffer.Delete()
		o.Buffer = nil
	}

	if o.Shader != nil {
		o.Shader.Delete()
		o.Shader = nil
	}

	if o.Mesh != nil {
		o.Mesh.Delete()
		o.Mesh = nil
	}

	if o.depthID != 0 {
		gl.DeleteRenderbuffers(1, &o.depthID)
		o.depthID = 0
	}

	if o.frameID != 0 {
		gl.DeleteFramebuffers(1, &o.frameID)
		o.frameID = 0
	}
}

// Update processes all events
func (o *Overlay) Update(ctx *context.Update) {

}

// Draw renders the current buffer to the screen
func (o *Overlay) Draw() {
	o.Shader.Bind()
	gl.UniformMatrix4fv(o.RenderCtx.GetShader().GetUniformLocation("uProjection"), 1, false, o.RenderCtx.GetProjectionPtr())

	gl.BindFramebuffer(gl.FRAMEBUFFER, o.frameID)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for _, c := range o.Components {
		c.Draw(&o.RenderCtx)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	gl.Uniform1i(o.Shader.GetUniformLocation("uTexture"), 0)
	gl.ActiveTexture(gl.TEXTURE0)
	o.Buffer.Bind()

	gl.Clear(gl.DEPTH_BUFFER_BIT)
	o.Mesh.Draw(&o.RenderCtx)

	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// AddComponent adds the given Component as a child node of the Overlay
func (o *Overlay) AddComponent(c Component) {
	o.Components = append(o.Components, c)
}
