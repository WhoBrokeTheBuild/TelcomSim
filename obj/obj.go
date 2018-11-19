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
	Material  *Material
}

type Material struct {
	Name                 string
	Ambient              mgl32.Vec3
	Diffuse              mgl32.Vec3
	Specular             mgl32.Vec3
	Shininess            float32
	Dissolve             float32
	AmbientMap           string
	DiffuseMap           string
	SpecularMap          string
	SpecularHighlightMap string
	BumpMap              string
	AlphaMap             string
	DisplacementMap      string
	ReflectionMap        string
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

func (rdr *reader) Read() ([]*Object, error) {
	rdr.filename = filepath.Clean(rdr.filename)
	dir := filepath.Dir(rdr.filename)

	objects := []*Object{}
	materials := map[string]*Material{}

	file, err := rdr.load(rdr.filename)
	if err != nil {
		return nil, err
	}

	if file[len(file)-1] != '\n' {
		file = append(file, '\n')
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
		bytes, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		line := string(bytes)

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
						&f[0][0], &f[0][1],
						&f[1][0], &f[1][1],
						&f[2][0], &f[2][1])
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

				if hasTxcd {
					o.TexCoords = append(o.TexCoords, txcds[f[i][1]-1])
				}

				if hasNorm {
					o.Normals = append(o.Normals, norms[f[i][2]-1])
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
			tmp, err := rdr.readMaterial(filepath.Join(dir, strings.TrimSpace(line[7:])))
			if err != nil {
				return nil, err
			}
			for k, v := range tmp {
				materials[k] = v
			}
		} else if strings.HasPrefix(line, "usemtl") {
			if o == nil {
				o = &Object{
					Name: "default",
				}
				objects = append(objects, o)
			}

			name := strings.TrimSpace(line[7:])
			if m, ok := materials[name]; ok {
				o.Material = m
			}
		}
	}

	return objects, nil
}

func (rdr *reader) readMaterial(filename string) (map[string]*Material, error) {
	filename = filepath.Clean(filename)
	dir := filepath.Dir(filename)

	materials := map[string]*Material{}

	file, err := rdr.load(filename)
	if err != nil {
		return nil, err
	}

	if file[len(file)-1] != '\n' {
		file = append(file, '\n')
	}

	buf := bytes.NewBuffer(file)

	var m *Material

	for {
		bytes, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		line := string(bytes)

		if len(line) < 2 || line[0] == '#' || line[0] == '\n' {
			continue
		}

		if strings.HasPrefix(line, "newmtl") { // newmtl
			name := strings.TrimSpace(line[7:])
			m = &Material{
				Name:     name,
				Ambient:  mgl32.Vec3{0, 0, 0},
				Diffuse:  mgl32.Vec3{0, 0, 0},
				Specular: mgl32.Vec3{0, 0, 0},
			}
			materials[name] = m
		} else if line[0] == 'K' {
			if line[1] == 'a' {
				if m != nil {
					// Ka
					fmt.Sscanf(line[3:], "%f %f %f", &m.Ambient[0], &m.Ambient[1], &m.Ambient[2])
				}
			} else if line[1] == 'd' {
				if m != nil {
					// Kd
					fmt.Sscanf(line[3:], "%f %f %f", &m.Diffuse[0], &m.Diffuse[1], &m.Diffuse[2])
				}
			} else if line[1] == 's' {
				if m != nil {
					// Ks
					fmt.Sscanf(line[3:], "%f %f %f", &m.Specular[0], &m.Specular[1], &m.Specular[2])
				}
			}
		} else if line[0] == 'N' && line[1] == 's' {
			// Ns
			fmt.Sscanf(line[3:], "%f", &m.Shininess)
		} else if strings.HasPrefix(line, "map_K") {
			if line[5] == 'a' {
				// map_Ka
				m.AmbientMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
			} else if line[5] == 'd' {
				// map_Kd
				m.DiffuseMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
			} else if line[5] == 's' {
				// map_Ks
				m.SpecularMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
			}
		} else if strings.HasPrefix(line, "map_Ns") {
			// map_Ns
			m.SpecularHighlightMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
		} else if strings.HasPrefix(line, "bump") {
			// bump
			m.BumpMap = filepath.Join(dir, strings.TrimSpace(line[5:]))
		} else if strings.HasPrefix(line, "map_bump") {
			// map_bump
			m.BumpMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
		} else if strings.HasPrefix(line, "disp") {
			// disp
			m.DisplacementMap = filepath.Join(dir, strings.TrimSpace(line[5:]))
		} else if strings.HasPrefix(line, "refl") {
			// refl
			m.ReflectionMap = filepath.Join(dir, strings.TrimSpace(line[5:]))
		} else if strings.HasPrefix(line, "map_d") {
			// map_d
			m.AlphaMap = filepath.Join(dir, strings.TrimSpace(line[6:]))
		}
	}

	return materials, nil
}
