package context

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Render is a context of view and shader data
type Render struct {
	View       mgl32.Mat4
	Projection mgl32.Mat4
}
