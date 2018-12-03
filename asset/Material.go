package asset

import (
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Material struct {
	Ambient  mgl32.Vec4
	Diffuse  mgl32.Vec4
	Specular mgl32.Vec4

	AmbientMap  *Texture
	DiffuseMap  *Texture
	SpecularMap *Texture
}

type MaterialData struct {
	Ambient  mgl32.Vec4
	Diffuse  mgl32.Vec4
	Specular mgl32.Vec4

	AmbientMap  string
	DiffuseMap  string
	SpecularMap string
}

const (
	// PositionAttrID is the attribute ID of _Position in GLSL
	PositionAttrID uint32 = 0
	// NormalAttrID is the attribute ID of _Normal in GLSL
	NormalAttrID uint32 = 1
	// TexCoordAttrID is the attribute ID of _TexCoord in GLSL
	TexCoordAttrID uint32 = 2
)

func NewMaterial(data *MaterialData) (*Material, error) {
	var err error
	m := &Material{
		Ambient:  data.Ambient,
		Diffuse:  data.Diffuse,
		Specular: data.Specular,
	}

	if data.AmbientMap != "" {
		m.AmbientMap, err = NewTextureFromFile(data.AmbientMap)
		if err != nil {
			return nil, err
		}
	}

	if data.DiffuseMap != "" {
		m.DiffuseMap, err = NewTextureFromFile(data.DiffuseMap)
		if err != nil {
			return nil, err
		}
	}

	if data.SpecularMap != "" {
		m.SpecularMap, err = NewTextureFromFile(data.SpecularMap)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

// Delete frees all resources owned by the Material
func (m *Material) Delete() {
	if m.AmbientMap != nil {
		m.AmbientMap.Delete()
		m.AmbientMap = nil
	}
	if m.DiffuseMap != nil {
		m.DiffuseMap.Delete()
		m.DiffuseMap = nil
	}
	if m.SpecularMap != nil {
		m.SpecularMap.Delete()
		m.SpecularMap = nil
	}
}

func (m *Material) Bind(s *Shader) {
	gl.Uniform1i(s.GetUniformLocation("uAmbientMap"), 0)
	if m.AmbientMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		m.AmbientMap.Bind()
		gl.Uniform4fv(s.GetUniformLocation("uAmbient"), 1, &[]float32{0, 0, 0, 0}[0])
	} else {
		gl.Uniform4fv(s.GetUniformLocation("uAmbient"), 1, &m.Ambient[0])
	}

	gl.Uniform1i(s.GetUniformLocation("uDiffuseMap"), 1)
	if m.DiffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		m.DiffuseMap.Bind()
		gl.Uniform4fv(s.GetUniformLocation("uDiffuse"), 1, &[]float32{0, 0, 0, 0}[0])
	} else {
		gl.Uniform4fv(s.GetUniformLocation("uDiffuse"), 1, &m.Diffuse[0])
	}

	gl.Uniform1i(s.GetUniformLocation("uSpecularMap"), 2)
	if m.SpecularMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		m.SpecularMap.Bind()
		gl.Uniform4fv(s.GetUniformLocation("uSpecular"), 1, &[]float32{0, 0, 0, 0}[0])
	} else {
		gl.Uniform4fv(s.GetUniformLocation("uSpecular"), 1, &m.Specular[0])
	}
}

func (m *Material) UnBind() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
