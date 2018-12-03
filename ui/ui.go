package ui

import (
	"github.com/WhoBrokeTheBuild/TelcomSim/asset"
	"github.com/go-gl/mathgl/mgl32"
)

func get2DMeshData(dst, src mgl32.Vec4) *asset.MeshData {
	return &asset.MeshData{
		Vertices: []mgl32.Vec3{
			mgl32.Vec3{dst[2], dst[1], 0},
			mgl32.Vec3{dst[2], dst[3], 0},
			mgl32.Vec3{dst[0], dst[3], 0},
			mgl32.Vec3{dst[2], dst[1], 0},
			mgl32.Vec3{dst[0], dst[3], 0},
			mgl32.Vec3{dst[0], dst[1], 0},
		},
		TexCoords: []mgl32.Vec2{
			mgl32.Vec2{src[2], src[3]},
			mgl32.Vec2{src[2], src[1]},
			mgl32.Vec2{src[0], src[1]},
			mgl32.Vec2{src[2], src[3]},
			mgl32.Vec2{src[0], src[1]},
			mgl32.Vec2{src[0], src[3]},
		},
	}
}

func new2DMesh(dst, src mgl32.Vec4) (*asset.Mesh, error) {
	mesh, err := asset.NewMesh(get2DMeshData(dst, src))
	if err != nil {
		return nil, err
	}
	return mesh, err
}

func update2DMesh(mesh *asset.Mesh, dst, src mgl32.Vec4) error {
	err := mesh.UpdateData(get2DMeshData(dst, src))
	if err != nil {
		return err
	}
	return err
}
