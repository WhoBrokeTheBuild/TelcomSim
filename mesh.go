package main

import (
	"C"

	"github.com/WhoBrokeTheBuild/TelcomSim/obj"
	gl "github.com/go-gl/gl/v3.1/gles2"
)

type Mesh struct {
	VAO   uint32
	VBO   uint32
	Count int32
}

const (
	InvalidID    uint32 = 0
	VertexAttrID uint32 = 0
)

func NewMesh(filename string) (*Mesh, error) {
	m := &Mesh{
		VAO: InvalidID,
		VBO: InvalidID,
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
	r := obj.NewReaderEx(filename, Asset)
	objs, err := r.Read()
	if err != nil {
		return err
	}

	o := objs[0]
	m.Count = int32(len(o.Vertices))

	gl.GenVertexArrays(1, &m.VAO)
	gl.BindVertexArray(m.VAO)

	gl.GenBuffers(1, &m.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(o.Vertices)*3*C.sizeof_float, gl.Ptr(o.Vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(VertexAttrID)
	gl.VertexAttribPointer(VertexAttrID, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	return nil
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, m.Count)
}
