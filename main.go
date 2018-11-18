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
	gl.ClearColor(0.39, 0.58, 0.93, 1.0)

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
	projectionUniform := gl.GetUniformLocation(s.ID, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(s.ID, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(s.ID, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

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
