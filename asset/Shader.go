package asset

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/WhoBrokeTheBuild/TelcomSim/data"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
	gl "github.com/go-gl/gl/v4.1-core/gl"
)

// Shader represents an OpenGL Shader Program
type Shader struct {
	ID       uint32
	Uniforms map[string]int32
}

var _versionString string

// NewShaderFromFiles returns a new Shader from the given files
func NewShaderFromFiles(filenames []string) (*Shader, error) {
	s := &Shader{
		ID:       InvalidID,
		Uniforms: map[string]int32{},
	}

	err := s.LoadFromFiles(filenames)
	if err != nil {
		s.Delete()
		return nil, err
	}

	return s, nil
}

// Delete frees all resources owned by the Shader
func (s *Shader) Delete() {
	if s.ID != InvalidID {
		gl.DeleteProgram(s.ID)
		s.ID = InvalidID
	}
}

// LoadFromFiles loads a shader from the given files
func (s *Shader) LoadFromFiles(filenames []string) error {
	s.Delete()

	shaders := make([]uint32, 0, len(filenames))
	for _, file := range filenames {
		id, err := compileShader(file)
		if err != nil {
			return err
		}
		shaders = append(shaders, id)
	}

	s.ID = gl.CreateProgram()
	for _, id := range shaders {
		gl.AttachShader(s.ID, id)
	}
	gl.LinkProgram(s.ID)

	var status int32
	gl.GetProgramiv(s.ID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(s.ID, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(s.ID, logLen, nil, gl.Str(log))

		s.Delete()
		return fmt.Errorf("Failed to link program: %v", log)
	}

	for _, id := range shaders {
		gl.DeleteShader(id)
	}

	s.cacheUniforms()

	return nil
}

// Bind calls glUseProgram with this Shader's ID
func (s *Shader) Bind() error {
	if s.ID == InvalidID {
		return fmt.Errorf("Failed to bind program: Not loaded")
	}

	gl.UseProgram(s.ID)
	return nil
}

// GetUniformLocation returns the uniform's location ID, or -1
func (s *Shader) GetUniformLocation(name string) int32 {
	if u, ok := s.Uniforms[name]; ok {
		return u
	}
	return -1
}

func (s *Shader) cacheUniforms() {
	var count int32

	var size int32
	var length int32
	var tp uint32
	buf := strings.Repeat("\x00", 256)

	gl.GetProgramiv(s.ID, gl.ACTIVE_UNIFORMS, &count)
	for i := int32(0); i < count; i++ {
		gl.GetActiveUniform(s.ID, uint32(i), int32(len(buf)), &length, &size, &tp, gl.Str(buf))

		// Force copy
		name := make([]byte, length)
		copy(name, []byte(buf[:length]))

		s.Uniforms[string(name)] = gl.GetUniformLocation(s.ID, gl.Str(string(name)+"\x00"))
	}
}

func getVersionString() string {
	if _versionString != "" {
		return _versionString
	}

	tmp := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	_versionString = "#version "
	for i := range tmp {
		if tmp[i] == ' ' {
			break
		}
		if unicode.IsDigit(rune(tmp[i])) {
			_versionString += string(tmp[i])
		}
	}
	_versionString += " core\n"
	return _versionString
}

func preProcessShader(code string) string {
	// Prepend `#version`
	code = getVersionString() + code

	// Append null-terminator (windows)
	code += "\x00"

	// Clean CRLF (windows)
	re := regexp.MustCompile(`\r`)
	code = re.ReplaceAllString(code, "")

	return code
}

func compileShader(filename string) (uint32, error) {
	filename = filepath.Clean(filename)

	t := getShaderType(filename)
	id := gl.CreateShader(t)

	log.Loadf("asset.Shader [%v]", filename)
	b, err := data.Asset(filename)
	if err != nil {
		return InvalidID, err
	}

	code := preProcessShader(string(b))

	re := regexp.MustCompile(`\r`)
	code = re.ReplaceAllString(code, "")

	ccode, free := gl.Strs(code)
	gl.ShaderSource(id, 1, ccode, nil)
	free()
	gl.CompileShader(id)

	var status int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(id, logLen, nil, gl.Str(log))

		return InvalidID, fmt.Errorf("Failed to compile [%v]: %v", filename, log)
	}

	return id, nil
}

func getShaderType(filename string) uint32 {
	if strings.HasSuffix(filename, ".vs.glsl") {
		return gl.VERTEX_SHADER
	}
	if strings.HasSuffix(filename, ".fs.glsl") {
		return gl.FRAGMENT_SHADER
	}
	return gl.INVALID_ENUM
}
