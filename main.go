package main

import (
	"fmt"
	"runtime"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  int = 1024
	WindowHeight int = 768
)

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

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

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

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(WindowWidth)/float32(WindowHeight), 0.1, 10.0)
	projectionUniform := s.GetUniformLocation("uProjection")
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	view := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	viewUniform := s.GetUniformLocation("uView")
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	model := mgl32.Ident4()
	modelUniform := s.GetUniformLocation("uModel")
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	color := fromRGB(249, 70, 8)
	colorUniform := s.GetUniformLocation("uColor")
	gl.Uniform4fv(colorUniform, 1, &color[0])

	m, err := NewMesh("assets/models/cube.obj")
	if err != nil {
		panic(err)
	}
	defer m.Delete()

	for !window.ShouldClose() {
		glfw.PollEvents()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		m.Draw()

		window.SwapBuffers()
	}
}
