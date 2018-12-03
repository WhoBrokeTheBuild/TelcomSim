package stbi

/*
#cgo CFLAGS: -DSTB_IMAGE_IMPLEMENTATION -DSTB_IMAGE_STATIC -DSTBI_NO_STDIO
#cgo !windows LDFLAGS: -lm
#include "stb_image.h"
*/
import "C"
