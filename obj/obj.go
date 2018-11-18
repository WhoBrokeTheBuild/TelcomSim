package obj

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

type Object struct {
	Name      string
	Vertices  []mgl32.Vec3
	Normals   []mgl32.Vec3
	TexCoords []mgl32.Vec2
}

type LoadFunc func(string) ([]byte, error)

type Reader interface {
	Read() ([]*Object, error)
}

func NewReader(filename string) Reader {
	return &reader{
		filename: filename,
		load:     ioutil.ReadFile,
	}
}

func NewReaderEx(filename string, load LoadFunc) Reader {
	return &reader{
		filename: filename,
		load:     load,
	}
}

type reader struct {
	filename string
	load     LoadFunc
}

func (r *reader) Read() ([]*Object, error) {
	objects := []*Object{}

	file, err := r.load(r.filename)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(file)

	var o *Object
	var x, y, z, u, v float32
	var f [3][3]int

	var hasNorm bool
	var hasTxcd bool

	verts := []mgl32.Vec3{}
	norms := []mgl32.Vec3{}
	txcds := []mgl32.Vec2{}

	for {
		b, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		line := string(b)

		if len(line) < 2 || line[0] == '#' || line[0] == '\n' {
			continue
		}

		if line[0] == 'v' {
			if line[1] == 'n' {
				fmt.Sscanf(line[3:], "%f %f %f", &x, &y, &z)
				norms = append(norms, mgl32.Vec3{x, y, z})
			} else if line[1] == 't' {
				fmt.Sscanf(line[3:], "%f %f", &u, &v)
				txcds = append(txcds, mgl32.Vec2{u, v})
			} else {
				fmt.Sscanf(line[2:], "%f %f %f", &x, &y, &z)
				verts = append(verts, mgl32.Vec3{x, y, z})
			}
		} else if line[0] == 'f' {
			if o == nil {
				o = &Object{
					Name: "default",
				}
				objects = append(objects, o)
			}

			hasNorm = false
			hasTxcd = false
			if strings.Contains(line, "//") {
				hasNorm = true
				fmt.Sscanf(line[2:], "%d//%d %d//%d %d//%d",
					&f[0][0], &f[0][2],
					&f[1][0], &f[1][2],
					&f[2][0], &f[2][2])
			} else {
				if strings.Count(line, "/") == 3 {
					hasTxcd = true
					fmt.Sscanf(line[2:], "%d/%d %d/%d %d/%d",
						&f[0][0], &f[0][2],
						&f[1][0], &f[1][2],
						&f[2][0], &f[2][2])
				} else {
					hasNorm = true
					hasTxcd = true
					fmt.Sscanf(line[2:], "%d/%d/%d %d/%d/%d %d/%d/%d",
						&f[0][0], &f[0][1], &f[0][2],
						&f[1][0], &f[1][1], &f[1][2],
						&f[2][0], &f[2][1], &f[2][2])
				}
			}
			// TODO: Handle `f %d %d %d`

			for i := 0; i < 3; i++ {
				if f[i][0] < 0 {
					f[i][0] += len(verts)
				}
				if f[i][1] < 0 {
					f[i][1] += len(txcds)
				}
				if f[i][2] < 0 {
					f[i][2] += len(norms)
				}

				o.Vertices = append(o.Vertices, verts[f[i][0]-1])

				if hasNorm {
					o.Normals = append(o.Normals, norms[f[i][2]-1])
				}

				if hasTxcd {
					o.TexCoords = append(o.TexCoords, txcds[f[i][1]-1])
				}
			}
		} else if line[0] == 'o' {
			o = &Object{
				Name:      strings.TrimSpace(line[2:]),
				Vertices:  []mgl32.Vec3{},
				Normals:   []mgl32.Vec3{},
				TexCoords: []mgl32.Vec2{},
			}
			objects = append(objects, o)
		} else if strings.HasPrefix(line, "mtllib") {
			fmt.Printf("Loading material [%s]\n", filepath.Join(
				filepath.Dir(r.filename),
				strings.TrimSpace(line[7:]),
			))
		} else if strings.HasPrefix(line, "usemtl") {
			if o == nil {
				o = &Object{
					Name: "default",
				}
				objects = append(objects, o)
			}
		}
	}

	return objects, nil
}
