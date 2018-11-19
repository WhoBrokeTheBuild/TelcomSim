package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	gl "github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	Name     string
	ID       uint32
	Uniforms map[string]int32
}

var _shaders = map[string]*Shader{}
var _versionString string

func NewShader(name string, sources []string) (*Shader, error) {
	s := &Shader{
		Name:     name,
		ID:       InvalidID,
		Uniforms: map[string]int32{},
	}

	if old, exists := _shaders[name]; exists {
		old.Delete()
	}

	err := s.Load(sources)
	if err != nil {
		s.Delete()
		return nil, err
	}

	_shaders[name] = s

	return s, nil
}

func (s *Shader) Delete() {
	if s.ID != InvalidID {
		gl.DeleteProgram(s.ID)
		s.ID = InvalidID
	}

	if _, exists := _shaders[s.Name]; exists {
		delete(_shaders, s.Name)
	}
}

func (s *Shader) Load(sources []string) error {
	s.Delete()

	shaders := make([]uint32, 0, len(sources))
	for _, src := range sources {
		id, err := compileShader(src)
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
		return fmt.Errorf("Failed to link program [%v]: %v", s.Name, log)
	}

	for _, id := range shaders {
		gl.DeleteShader(id)
	}

	s.cacheUniforms()

	return nil
}

func (s *Shader) Bind() error {
	if s.ID == InvalidID {
		return fmt.Errorf("Failed to bind program [%v]: Not loaded", s.Name)
	}

	gl.UseProgram(s.ID)
	return nil
}

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

	Loadf("Shader [%v]", filename)
	b, err := LoadAsset(filename)
	if err != nil {
		return InvalidID, fmt.Errorf("Failed to load file [%v]", filename)
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
