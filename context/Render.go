package context

import (
	"github.com/WhoBrokeTheBuild/TelcomSim/asset"
	"github.com/go-gl/mathgl/mgl32"
)

// Render is a context of view and shader data
type Render struct {
	Projection mgl32.Mat4
	View       mgl32.Mat4
	Shader     *asset.Shader
}

// GetProjectionPtr returns the projection matrix as a *float32
func (r *Render) GetProjectionPtr() *float32 {
	return &r.Projection[0]
}

// GetViewPtr returns the view matrix as a *float32
func (r *Render) GetViewPtr() *float32 {
	return &r.View[0]
}

// GetShader returns the Shader
func (r *Render) GetShader() *asset.Shader {
	return r.Shader
}
