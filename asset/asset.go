package asset

const (
	// InvalidID is an invalid OpenGL ID
	InvalidID uint32 = 0
)

// Used to decouple `asset` from `context`
type renderContext interface {
	GetProjectionPtr() *float32
	GetViewPtr() *float32
	GetShader() *Shader
}
