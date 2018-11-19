package main

import (
	"C"

	"github.com/WhoBrokeTheBuild/TelcomSim/obj"
	gl "github.com/go-gl/gl/v4.1-core/gl"
)
import (
	"fmt"
	"path/filepath"

	"github.com/go-gl/mathgl/mgl32"
)

type Mesh struct {
	VAO   uint32
	VBO   uint32
	Count int32
	Model mgl32.Mat4

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	AmbientMap  *Texture
	DiffuseMap  *Texture
	SpecularMap *Texture
}

const (
	PositionAttrID uint32 = 0
	NormalAttrID   uint32 = 1
	TexCoordAttrID uint32 = 2
)

func NewMesh(filename string) (*Mesh, error) {
	m := &Mesh{
		VAO:   InvalidID,
		VBO:   InvalidID,
		Model: mgl32.Ident4(),
	}

	err := m.Load(filename)
	if err != nil {
		m.Delete()
		return nil, err
	}

	return m, nil
}

func (m *Mesh) Delete() {
	if m.VBO != InvalidID {
		gl.DeleteBuffers(1, &m.VBO)
		m.VBO = InvalidID
	}
	if m.VAO != InvalidID {
		gl.DeleteVertexArrays(1, &m.VAO)
		m.VAO = InvalidID
	}
}

func (m *Mesh) Load(filename string) error {
	filename = filepath.Clean(filename)
	m.Delete()

	Loadf("Mesh [%v]", filename)
	r := obj.NewReaderEx(filename, LoadAsset)
	objs, err := r.Read()
	if err != nil {
		return err
	}

	if len(objs) == 0 {
		return fmt.Errorf("No objects loaded from [%v]", filename)
	}

	o := objs[0]

	if o.Material != nil {
		m.Ambient = o.Material.Ambient
		m.Diffuse = o.Material.Diffuse
		m.Specular = o.Material.Specular

		if o.Material.AmbientMap != "" {
			m.AmbientMap, err = NewTexture(o.Material.AmbientMap)
			if err != nil {
				return err
			}
		}

		if o.Material.DiffuseMap != "" {
			m.DiffuseMap, err = NewTexture(o.Material.DiffuseMap)
			if err != nil {
				return err
			}
		}

		if o.Material.SpecularMap != "" {
			m.SpecularMap, err = NewTexture(o.Material.SpecularMap)
			if err != nil {
				return err
			}
		}
	}

	m.Count = int32(len(o.Vertices))
	hasNorms := len(o.Normals) > 0
	hasTxcds := len(o.TexCoords) > 0

	buf := make([]float32, 0, (len(o.Vertices)*3)+(len(o.Normals)*3)+(len(o.TexCoords)*2))
	for i := range o.Vertices {
		buf = append(buf, o.Vertices[i][0], o.Vertices[i][1], o.Vertices[i][2])
		if hasNorms {
			buf = append(buf, o.Normals[i][0], o.Normals[i][1], o.Normals[i][2])
		} else {
			buf = append(buf, 0.0, 0.0, 0.0)
		}
		if hasTxcds {
			buf = append(buf, o.TexCoords[i][0], o.TexCoords[i][1])
		} else {
			buf = append(buf, 0.0, 0.0)
		}
	}

	F := C.sizeof_float
	stride := int32(8 * F)

	gl.GenVertexArrays(1, &m.VAO)
	gl.BindVertexArray(m.VAO)

	gl.GenBuffers(1, &m.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(buf)*F, gl.Ptr(buf), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(PositionAttrID)
	gl.VertexAttribPointer(PositionAttrID, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(NormalAttrID)
	gl.VertexAttribPointer(NormalAttrID, 3, gl.FLOAT, false, stride, gl.PtrOffset(3*F))

	gl.EnableVertexAttribArray(TexCoordAttrID)
	gl.VertexAttribPointer(TexCoordAttrID, 2, gl.FLOAT, false, stride, gl.PtrOffset(6*F))

	return nil
}

func (m *Mesh) Draw(ctx *RenderContext) {
	gl.UniformMatrix4fv(ctx.Shader.GetUniformLocation("uProjection"), 1, false, &ctx.Projection[0])
	gl.UniformMatrix4fv(ctx.Shader.GetUniformLocation("uView"), 1, false, &ctx.View[0])
	gl.UniformMatrix4fv(ctx.Shader.GetUniformLocation("uModel"), 1, false, &m.Model[0])

	gl.Uniform3fv(ctx.Shader.GetUniformLocation("uAmbient"), 1, &m.Ambient[0])
	gl.Uniform3fv(ctx.Shader.GetUniformLocation("uDiffuse"), 1, &m.Diffuse[0])
	gl.Uniform3fv(ctx.Shader.GetUniformLocation("uSpecular"), 1, &m.Specular[0])

	if m.AmbientMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		m.AmbientMap.Bind()
		gl.Uniform1i(ctx.Shader.GetUniformLocation("uAmbientMap"), 0)
	}

	if m.DiffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		m.DiffuseMap.Bind()
		gl.Uniform1i(ctx.Shader.GetUniformLocation("uDiffuseMap"), 1)
	}

	if m.SpecularMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		m.SpecularMap.Bind()
		gl.Uniform1i(ctx.Shader.GetUniformLocation("uSpecularMap"), 2)
	}

	gl.BindVertexArray(m.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, m.Count)
}
