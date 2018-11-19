package main

import (
	"runtime"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  int = 1024
	WindowHeight int = 768

	InvalidID uint32 = 0
)

type UpdateContext struct {
	DeltaTime   float32
	ElapsedTime float64
	TotalTime   float64
}

type RenderContext struct {
	View       mgl32.Mat4
	Projection mgl32.Mat4
	Shader     *Shader
}

func init() {
	runtime.LockOSThread()
}

func fromRGB(r, g, b int) mgl32.Vec4 {
	return mgl32.Vec4{
		float32(r) / float32(255),
		float32(g) / float32(255),
		float32(b) / float32(255),
		1.0,
	}
}

func main() {
	var err error

	err = glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Telcom Simulator", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	Infof("OpenGL Version %s", gl.GoStr(gl.GetString(gl.VERSION)))
	Infof("GLSL Version %s", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	Infof("OpenGL Vendor %s", gl.GoStr(gl.GetString(gl.VENDOR)))
	Infof("OpenGL Renderer %s", gl.GoStr(gl.GetString(gl.RENDERER)))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	cc := fromRGB(16, 163, 160)
	gl.ClearColor(cc[0], cc[1], cc[2], 1.0)

	s, err := NewShader("default", []string{
		"assets/shaders/default.vs.glsl",
		"assets/shaders/default.fs.glsl",
	})
	if err != nil {
		panic(err)
	}
	defer s.Delete()

	s.Bind()

	light := mgl32.Vec3{3, 3, 0}
	camera := mgl32.Vec3{3, 3, 3}

	gl.Uniform3fv(s.GetUniformLocation("uLight"), 1, &light[0])
	gl.Uniform3fv(s.GetUniformLocation("uCamera"), 1, &camera[0])

	//m, err := NewMesh("assets/models/crate/crate.obj")
	m, err := NewMesh("assets/models/earth/earth.obj")
	if err != nil {
		panic(err)
	}
	defer m.Delete()

	updateCtx := &UpdateContext{0, 0, 0}
	renderCtx := &RenderContext{
		View:       mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}),
		Projection: mgl32.Perspective(mgl32.DegToRad(45.0), float32(WindowWidth)/float32(WindowHeight), 0.1, 10.0),
		Shader:     s,
	}

	rotation := 0.0
	previousTime := glfw.GetTime()
	for !window.ShouldClose() {
		glfw.PollEvents()

		time := glfw.GetTime()
		updateCtx.ElapsedTime = time - previousTime
		previousTime = time

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		rotation += updateCtx.ElapsedTime
		m.Model = mgl32.HomogRotate3D(float32(rotation), mgl32.Vec3{0, 1, 0})

		m.Draw(renderCtx)

		window.SwapBuffers()
	}
}
