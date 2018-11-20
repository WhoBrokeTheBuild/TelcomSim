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

type meshGroup struct {
	VAO   uint32
	VBO   uint32
	Count int32

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	AmbientMap  *Texture
	DiffuseMap  *Texture
	SpecularMap *Texture
}

type Mesh struct {
	Model  mgl32.Mat4
	groups []*meshGroup
}

type MeshData struct {
	Vertices  []mgl32.Vec3
	Normals   []mgl32.Vec3
	TexCoords []mgl32.Vec2

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	AmbientMap  string
	DiffuseMap  string
	SpecularMap string
}

const (
	PositionAttrID uint32 = 0
	NormalAttrID   uint32 = 1
	TexCoordAttrID uint32 = 2
)

func NewMeshFromData(data []MeshData) (*Mesh, error) {
	m := &Mesh{
		Model:  mgl32.Ident4(),
		groups: []*meshGroup{},
	}

	err := m.LoadFromData(data)
	if err != nil {
		m.Delete()
		return nil, err
	}

	return m, err
}

func NewMeshFromFile(filename string) (*Mesh, error) {
	m := &Mesh{
		Model:  mgl32.Ident4(),
		groups: []*meshGroup{},
	}

	err := m.LoadFromFile(filename)
	if err != nil {
		m.Delete()
		return nil, err
	}

	return m, nil
}

func (m *Mesh) Delete() {
	for _, g := range m.groups {
		if g.VBO != InvalidID {
			gl.DeleteBuffers(1, &g.VBO)
			g.VBO = InvalidID
		}
		if g.VAO != InvalidID {
			gl.DeleteVertexArrays(1, &g.VAO)
			g.VAO = InvalidID
		}
	}
}

func (m *Mesh) LoadFromFile(filename string) error {
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

	data := make([]MeshData, 0, len(objs))

	for _, o := range objs {
		data = append(data, MeshData{
			Vertices: o.Vertices,
			Normals: o.Normals,
			TexCoords: o.TexCoords,
			Ambient: o.Material.Ambient,
			Diffuse: o.Material.Diffuse,
			Specular: o.Material.Specular,
			AmbientMap: o.Material.AmbientMap,
			DiffuseMap: o.Material.DiffuseMap,
			SpecularMap: o.Material.SpecularMap,
		})
	}

	err = m.LoadFromData(data)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mesh) LoadFromData(data []MeshData) error {
	var err error 
	
	for _, d := range data {
		g := &meshGroup{
			VAO: InvalidID,
			VBO: InvalidID,
		}
		m.groups = append(m.groups, g)

			g.Ambient = d.Ambient
			g.Diffuse = d.Diffuse
			g.Specular = d.Specular

			if d.AmbientMap != "" {
				g.AmbientMap, err = NewTexture(d.AmbientMap)
				if err != nil {
					return err
				}
			}

			if d.DiffuseMap != "" {
				g.DiffuseMap, err = NewTexture(d.DiffuseMap)
				if err != nil {
					return err
				}
			}

			if d.SpecularMap != "" {
				g.SpecularMap, err = NewTexture(d.SpecularMap)
				if err != nil {
					return err
				}
			}

		g.Count = int32(len(d.Vertices))
		hasNorms := len(d.Normals) > 0
		hasTxcds := len(d.TexCoords) > 0

		buf := make([]float32, 0, (len(d.Vertices)*3)+(len(d.Normals)*3)+(len(d.TexCoords)*2))
		for i := range d.Vertices {
			buf = append(buf, d.Vertices[i][0], d.Vertices[i][1], d.Vertices[i][2])
			if hasNorms {
				buf = append(buf, d.Normals[i][0], d.Normals[i][1], d.Normals[i][2])
			} else {
				buf = append(buf, 0.0, 0.0, 0.0)
			}
			if hasTxcds {
				buf = append(buf, d.TexCoords[i][0], d.TexCoords[i][1])
			} else {
				buf = append(buf, 0.0, 0.0)
			}
		}

		F := C.sizeof_float
		stride := int32(8 * F)

		gl.GenVertexArrays(1, &g.VAO)
		gl.BindVertexArray(g.VAO)

		gl.GenBuffers(1, &g.VBO)
		gl.BindBuffer(gl.ARRAY_BUFFER, g.VBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(buf)*F, gl.Ptr(buf), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(PositionAttrID)
		gl.VertexAttribPointer(PositionAttrID, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))

		gl.EnableVertexAttribArray(NormalAttrID)
		gl.VertexAttribPointer(NormalAttrID, 3, gl.FLOAT, false, stride, gl.PtrOffset(3*F))

		gl.EnableVertexAttribArray(TexCoordAttrID)
		gl.VertexAttribPointer(TexCoordAttrID, 2, gl.FLOAT, false, stride, gl.PtrOffset(6*F))
	}

	return nil
}

func (m *Mesh) Draw(ctx *RenderContext) {
	ctx.Shader.Bind()
	
	gl.UniformMatrix4fv(ctx.Shader.GetUniformLocation("uProjection"), 1, false, &ctx.Projection[0])
	gl.UniformMatrix4fv(ctx.Shader.GetUniformLocation("uView"), 1, false, &ctx.View[0])
	gl.UniformMatrix4fv(ctx.Shader.GetUniformLocation("uModel"), 1, false, &m.Model[0])

	for _, g := range m.groups {
		gl.Uniform3fv(ctx.Shader.GetUniformLocation("uAmbient"), 1, &g.Ambient[0])
		gl.Uniform3fv(ctx.Shader.GetUniformLocation("uDiffuse"), 1, &g.Diffuse[0])
		gl.Uniform3fv(ctx.Shader.GetUniformLocation("uSpecular"), 1, &g.Specular[0])

		if g.AmbientMap != nil {
			gl.ActiveTexture(gl.TEXTURE0)
			g.AmbientMap.Bind()
			gl.Uniform1i(ctx.Shader.GetUniformLocation("uAmbientMap"), 0)
		}

		if g.DiffuseMap != nil {
			gl.ActiveTexture(gl.TEXTURE1)
			g.DiffuseMap.Bind()
			gl.Uniform1i(ctx.Shader.GetUniformLocation("uDiffuseMap"), 1)
		}

		if g.SpecularMap != nil {
			gl.ActiveTexture(gl.TEXTURE2)
			g.SpecularMap.Bind()
			gl.Uniform1i(ctx.Shader.GetUniformLocation("uSpecularMap"), 2)
		}

		gl.BindVertexArray(g.VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, g.Count)
	}
}
