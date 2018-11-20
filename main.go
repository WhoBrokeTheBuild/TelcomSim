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

var defaultShader *Shader
var uiShader *Shader

var uiMesh *Mesh

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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	cc := fromRGB(16, 163, 160)
	gl.ClearColor(cc[0], cc[1], cc[2], 1.0)

	defaultShader, err = NewShader("default", []string{
		"assets/shaders/default.vs.glsl",
		"assets/shaders/default.fs.glsl",
	})
	if err != nil {
		panic(err)
	}
	defer defaultShader.Delete()

	uiShader, err = NewShader("ui", []string{
		"assets/shaders/ui.vs.glsl",
		"assets/shaders/ui.fs.glsl",
	})
	if err != nil {
		panic(err)
	}
	defer uiShader.Delete()

	uiMeshData := MeshData{
		Vertices: []mgl32.Vec3{
			mgl32.Vec3{float32(WindowWidth), 0, 0},
			mgl32.Vec3{float32(WindowWidth), float32(WindowHeight), 0},
			mgl32.Vec3{0, float32(WindowHeight), 0},
			mgl32.Vec3{float32(WindowWidth), 0, 0},
			mgl32.Vec3{0, float32(WindowHeight), 0},
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
		DiffuseMap: "assets/ui-test.png",
	}
	uiMesh, err = NewMeshFromData([]MeshData{uiMeshData})
	if err != nil {
		panic(err)
	}
	defer uiMesh.Delete()

	defaultShader.Bind()

	camera := mgl32.Vec3{3, 3, 3}
	gl.Uniform3fv(defaultShader.GetUniformLocation("uCamera"), 1, &camera[0])

	light := mgl32.Vec3{0, 5, 5}
	gl.Uniform3fv(defaultShader.GetUniformLocation("uLight"), 1, &light[0])

	m, err := NewMeshFromFile("assets/models/crate/crate.obj")
	if err != nil {
		panic(err)
	}
	defer m.Delete()

	aspect := float32(WindowWidth) / float32(WindowHeight)

	updateCtx := &UpdateContext{0, 0, 0}
	defaultRenderCtx := &RenderContext{
		View:       mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 1, 0}),
		Projection: mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0),
		Shader:     defaultShader,
	}
	uiRenderCtx := &RenderContext{
		Projection: mgl32.Ortho2D(0, float32(WindowWidth), float32(WindowHeight), 0),
		Shader:     uiShader,
	}

	rotation := 0.0
	previousTime := glfw.GetTime()
	for !window.ShouldClose() {
		glfw.PollEvents()

		if window.GetKey(glfw.KeyF2) == glfw.Press {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else if window.GetKey(glfw.KeyF2) == glfw.Release {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}

		time := glfw.GetTime()
		updateCtx.ElapsedTime = time - previousTime
		previousTime = time

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		rotation += updateCtx.ElapsedTime
		m.Model = mgl32.HomogRotate3D(float32(rotation), mgl32.Vec3{0, 1, 0})
		m.Draw(defaultRenderCtx)

		gl.Clear(gl.DEPTH_BUFFER_BIT)

		uiMesh.Draw(uiRenderCtx)

		window.SwapBuffers()
	}
}
