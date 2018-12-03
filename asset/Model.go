package asset

import (
	"fmt"
	"path/filepath"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/WhoBrokeTheBuild/TelcomSim/data"
	"github.com/WhoBrokeTheBuild/TelcomSim/log"
	"github.com/WhoBrokeTheBuild/TelcomSim/obj"
	"github.com/go-gl/mathgl/mgl32"
)

type Model struct {
	Transform mgl32.Mat4
	Meshes    []*Mesh
}

// NewModelFromFile returns a new Model from the given file
func NewModelFromFile(filename string) (*Model, error) {
	m := &Model{
		Transform: mgl32.Ident4(),
		Meshes:    []*Mesh{},
	}

	err := m.LoadFromFile(filename)
	if err != nil {
		m.Delete()
		return nil, err
	}

	return m, nil
}

// Delete frees all resources owned by the Model
func (m *Model) Delete() {
	for _, mesh := range m.Meshes {
		mesh.Delete()
	}
	m.Meshes = []*Mesh{}
}

// LoadFromFile loads a mesh from a given file
func (m *Model) LoadFromFile(filename string) error {
	filename = filepath.Clean(filename)
	m.Delete()

	log.Loadf("asset.Model [%v]", filename)
	r := obj.NewReaderEx(filename, obj.LoadFunc(data.Asset))
	objs, err := r.Read()
	if err != nil {
		return err
	}

	if len(objs) == 0 {
		return fmt.Errorf("No objects loaded from [%v]", filename)
	}

	for _, o := range objs {
		mat, err := NewMaterial(&MaterialData{
			Ambient:     mgl32.Vec4{o.Material.Ambient[0], o.Material.Ambient[1], o.Material.Ambient[2], 1},
			Diffuse:     mgl32.Vec4{o.Material.Diffuse[0], o.Material.Diffuse[1], o.Material.Diffuse[2], 1},
			Specular:    mgl32.Vec4{o.Material.Specular[0], o.Material.Specular[1], o.Material.Specular[2], 1},
			AmbientMap:  o.Material.AmbientMap,
			DiffuseMap:  o.Material.DiffuseMap,
			SpecularMap: o.Material.SpecularMap,
		})
		if err != nil {
			return err
		}
		mesh, err := NewMesh(&MeshData{
			Material:  mat,
			Vertices:  o.Vertices,
			Normals:   o.Normals,
			TexCoords: o.TexCoords,
		})
		if err != nil {
			return err
		}
		m.Meshes = append(m.Meshes, mesh)
	}

	return nil
}

// Draw renders a Model to the screen
func (m *Model) Draw(ctx renderContext) {
	ctx.GetShader().Bind()

	gl.UniformMatrix4fv(ctx.GetShader().GetUniformLocation("uProjection"), 1, false, ctx.GetProjectionPtr())
	gl.UniformMatrix4fv(ctx.GetShader().GetUniformLocation("uView"), 1, false, ctx.GetViewPtr())
	gl.UniformMatrix4fv(ctx.GetShader().GetUniformLocation("uModel"), 1, false, &m.Transform[0])

	for _, mesh := range m.Meshes {
		mesh.Draw(ctx)
	}
}
