package main

import (
	"fmt"
	"strings"

	gl "github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	Name string
	ID   uint32
}

const (
	InvalidShaderID  uint32 = 0
	InvalidProgramID uint32 = 0
)

var _shaders = map[string]*Shader{}

func NewShader(name string, sources []string) (*Shader, error) {
	s := &Shader{
		Name: name,
		ID:   InvalidProgramID,
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
	if s.ID != InvalidProgramID {
		gl.DeleteProgram(s.ID)
		s.ID = InvalidProgramID
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

	return nil
}

func (s *Shader) Bind() error {
	if s.ID == InvalidProgramID {
		return fmt.Errorf("Failed to bind program [%v]: Not loaded", s.Name)
	}

	gl.UseProgram(s.ID)
	return nil
}

func (s *Shader) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
}

func compileShader(filename string) (uint32, error) {
	t := getShaderType(filename)
	id := gl.CreateShader(t)

	b, err := Asset(filename)
	if err != nil {
		return InvalidShaderID, fmt.Errorf("Failed to load file [%v]", filename)
	}

	code := string(b)

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

		return InvalidShaderID, fmt.Errorf("Failed to compile [%v]: %v", filename, log)
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
