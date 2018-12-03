package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"runtime"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/TelcomSim/asset"
	"github.com/WhoBrokeTheBuild/TelcomSim/context"
	"github.com/WhoBrokeTheBuild/TelcomSim/data"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
	"github.com/WhoBrokeTheBuild/TelcomSim/ui"
)

const (
	windowWidth  int = 1024
	windowHeight int = 768
)

func init() {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
}

func fromRGB(r, g, b int) mgl32.Vec4 {
	return mgl32.Vec4{
		float32(r) / float32(255),
		float32(g) / float32(255),
		float32(b) / float32(255),
		1.0,
	}
}

var defaultShader *asset.Shader

var hud *ui.Overlay
var fps *ui.Text

func main() {
	var err error

	log.Infof("CPU Cores: %v", runtime.NumCPU()-1)

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
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Telcom Simulator", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	icons := []image.Image{}
	iconFiles := []string{"icons/icon_256.png", "icons/icon_128.png", "icons/icon_64.png", "icons/icon_32.png"}
	for _, file := range iconFiles {
		b, err := data.Asset(file)
		if err != nil {
			continue
		}
		image, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			continue
		}
		icons = append(icons, image)
	}
	window.SetIcon(icons)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	log.Infof("OpenGL Version: [%s]", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Infof("GLSL Version: [%s]", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	log.Infof("OpenGL Vendor: [%s]", gl.GoStr(gl.GetString(gl.VENDOR)))
	log.Infof("OpenGL Renderer: [%s]", gl.GoStr(gl.GetString(gl.RENDERER)))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	//music, err := asset.NewSoundFromFile("music/bgm.mp3")
	//if err != nil {
	//	panic(err)
	//}

	hud, err = ui.NewOverlay(mgl32.Vec2{float32(windowWidth), float32(windowHeight)})
	if err != nil {
		panic(err)
	}
	defer hud.Delete()

	initUI()

	defaultShader, err = asset.NewShaderFromFiles([]string{
		"shaders/default.vs.glsl",
		"shaders/default.fs.glsl",
	})
	if err != nil {
		panic(err)
	}
	defer defaultShader.Delete()

	defaultShader.Bind()

	camera := mgl32.Vec3{3, 3, 3}
	gl.Uniform3fv(defaultShader.GetUniformLocation("uCamera"), 1, &camera[0])

	light := mgl32.Vec3{3, 3, 3}
	gl.Uniform3fv(defaultShader.GetUniformLocation("uLight"), 1, &light[0])

	m, err := asset.NewModelFromFile("models/crate/crate.obj")
	if err != nil {
		panic(err)
	}
	defer m.Delete()

	aspect := float32(windowWidth) / float32(windowHeight)

	updateCtx := &context.Update{}
	renderCtx := &context.Render{
		View:       mgl32.LookAtV(mgl32.Vec3{2, 2, 2}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}),
		Projection: mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0),
		Shader:     defaultShader,
	}

	rotation := 0.0
	update := func(ctx *context.Update) {
		if window.GetKey(glfw.KeyF2) == glfw.Press {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else if window.GetKey(glfw.KeyF2) == glfw.Release {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}

		hud.Update(updateCtx)
		rotation += ctx.ElapsedTime

		m.Transform = mgl32.HomogRotate3D(float32(rotation), mgl32.Vec3{0, 1, 0})
	}

	render := func(ctx *context.Render) {
		m.Draw(renderCtx)
		hud.Draw()
	}

	//music.Play()

	const (
		frameDelay = 1.0 / 60.0
		fpsDelay   = 1.0
	)
	var (
		frameElap  = 0.0
		frameCount = 0
		fpsElap    = 0.0
	)

	prev := glfw.GetTime()
	for !window.ShouldClose() {
		time := glfw.GetTime()
		elapsed := time - prev
		prev = time

		frameElap += elapsed
		fpsElap += elapsed
		updateCtx.DeltaTime = 0.0
		updateCtx.ElapsedTime = elapsed

		if fpsElap >= fpsDelay {
			if fps != nil {
				if frameCount < 30 {
					fps.Color = color.RGBA{255, 0, 0, 255}
				} else if frameCount < 60 {
					fps.Color = color.RGBA{255, 255, 0, 255}
				} else {
					fps.Color = color.RGBA{0, 255, 0, 255}
				}
				fps.SetText(fmt.Sprintf("FPS %d", frameCount))
			}
			fpsElap = 0.0
			frameCount = 0
		}

		glfw.PollEvents()
		update(updateCtx)

		if frameElap >= frameDelay {
			frameCount++
			frameElap = 0.0

			cc := fromRGB(16, 163, 160)
			gl.ClearColor(cc[0], cc[1], cc[2], 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			render(renderCtx)

			window.SwapBuffers()
		}
	}
}

func initUI() {
	hud.AddComponent(ui.NewImageFromFile("ui/menubar.png"))

	fps = ui.NewText("FPS 00", "ui/default.ttf", 18.0, color.White)
	fps.SetPosition(mgl32.Vec2{float32(windowWidth) - 60, 5})
	hud.AddComponent(fps)

	menu := ui.NewText("File  Edit  Window", "ui/default.ttf", 18.0, color.White)
	menu.SetPosition(mgl32.Vec2{10, 5})
	hud.AddComponent(menu)

	//box := ui.NewImageFromFile("models/crate/crate.png")
	//box.SetPosition(mgl32.Vec2{100, 100})
	//hud.AddComponent(box)
}
