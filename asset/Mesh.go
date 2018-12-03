package asset

import (
	"C"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Mesh represents a set of OpenGL Vertex Array Objects and Material data
type Mesh struct {
	Material *Material
	VAO      uint32
	VBO      uint32
	Size     int
	Count    int32
}

// MeshData is the intermediate data format for loading Meshes from Memory
type MeshData struct {
	Material *Material

	Vertices  []mgl32.Vec3
	Normals   []mgl32.Vec3
	TexCoords []mgl32.Vec2
}

// NewMesh returns a new Mesh from the given MeshData
func NewMesh(data *MeshData) (*Mesh, error) {
	m := &Mesh{
		Material: data.Material,
	}

	err := m.LoadFromData(data)
	if err != nil {
		m.Delete()
		return nil, err
	}

	return m, err
}

// Delete frees all resources owned by the Mesh
func (m *Mesh) Delete() {
	if m.Material != nil {
		m.Material.Delete()
	}
	if m.VBO != InvalidID {
		gl.DeleteBuffers(1, &m.VBO)
		m.VBO = InvalidID
	}
	if m.VAO != InvalidID {
		gl.DeleteVertexArrays(1, &m.VAO)
		m.VAO = InvalidID
	}
}

// LoadFromData loads a mesh from an array of MeshData
func (m *Mesh) LoadFromData(data *MeshData) error {
	const F = C.sizeof_float

	m.Count = int32(len(data.Vertices))
	hasNorms := len(data.Normals) > 0
	hasTxcds := len(data.TexCoords) > 0

	buf := make([]float32, 0, (len(data.Vertices)*3)+(len(data.Normals)*3)+(len(data.TexCoords)*2))
	for i := range data.Vertices {
		buf = append(buf, data.Vertices[i][0], data.Vertices[i][1], data.Vertices[i][2])
		if hasNorms {
			buf = append(buf, data.Normals[i][0], data.Normals[i][1], data.Normals[i][2])
		}
		if hasTxcds {
			buf = append(buf, data.TexCoords[i][0], data.TexCoords[i][1])
		}
	}

	m.Size = len(buf)

	stride := int32(3 * F)
	if hasNorms {
		stride += int32(3 * F)
	}
	if hasTxcds {
		stride += int32(2 * F)
	}

	offset := 0

	gl.GenVertexArrays(1, &m.VAO)
	gl.BindVertexArray(m.VAO)

	gl.GenBuffers(1, &m.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(buf)*F, gl.Ptr(buf), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(PositionAttrID)
	gl.VertexAttribPointer(PositionAttrID, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	offset += 3 * F

	if hasNorms {
		gl.EnableVertexAttribArray(NormalAttrID)
		gl.VertexAttribPointer(NormalAttrID, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
		offset += 3 * F
	}

	if hasTxcds {
		gl.EnableVertexAttribArray(TexCoordAttrID)
		gl.VertexAttribPointer(TexCoordAttrID, 2, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return nil
}

// UpdateData sets the data in the existing buffer
func (m *Mesh) UpdateData(data *MeshData) error {
	const F = C.sizeof_float

	m.Count = int32(len(data.Vertices))
	hasNorms := len(data.Normals) > 0
	hasTxcds := len(data.TexCoords) > 0

	buf := make([]float32, 0, (len(data.Vertices)*3)+(len(data.Normals)*3)+(len(data.TexCoords)*2))
	for i := range data.Vertices {
		buf = append(buf, data.Vertices[i][0], data.Vertices[i][1], data.Vertices[i][2])
		if hasNorms {
			buf = append(buf, data.Normals[i][0], data.Normals[i][1], data.Normals[i][2])
		}
		if hasTxcds {
			buf = append(buf, data.TexCoords[i][0], data.TexCoords[i][1])
		}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)

	if m.Size == len(buf) {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(buf)*F, gl.Ptr(buf))
	} else {
		m.Size = len(buf)
		gl.BufferData(gl.ARRAY_BUFFER, len(buf)*F, gl.Ptr(buf), gl.STATIC_DRAW)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

// Draw renders a Mesh to the screen
func (m *Mesh) Draw(ctx renderContext) {
	if m.Material != nil {
		m.Material.Bind(ctx.GetShader())
	}

	gl.BindVertexArray(m.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, m.Count)

	if m.Material != nil {
		m.Material.UnBind()
	}
}
