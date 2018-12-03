package stbi

/*
#include "stb_image.h"
*/
import "C"
import "unsafe"

const (
	// True = 1
	True C.int = 1
	// False = 0
	False C.int = 0
	// Null = 0
	Null C.int = 0

	// Default = STBI_default
	Default C.int = C.STBI_default
	// Grey = STBI_grey
	Grey C.int = C.STBI_grey
	// GreyAlpha = STBI_grey_alpha
	GreyAlpha C.int = C.STBI_grey_alpha
	// RGB = STBI_rgb
	RGB C.int = C.STBI_rgb
	// RGBAlpha = STBI_rgb_alpha
	RGBAlpha C.int = C.STBI_rgb_alpha
)

// SetFlipVerticallyOnLoad = stbi_set_flip_vertically_on_load
func SetFlipVerticallyOnLoad(f C.int) {
	C.stbi_set_flip_vertically_on_load(f)
}

// Load = stbi_load
//func Load(filename string, desiredChannels C.int) (*C.uchar, int, int, C.int) {
//	var width, height, channels C.int
//	cstr := C.CString(filename)
//	defer C.free(unsafe.Pointer(cstr))
//
//	return C.stbi_load(cstr, &width, &height, &channels, desiredChannels),
//		int(width), int(height), channels
//}

// LoadFromMemory = stbi_load_from_memory
func LoadFromMemory(buffer []byte, desiredChannels C.int) (*C.uchar, int, int, C.int) {
	var width, height, channels C.int
	cbuf := C.CBytes(buffer)
	defer C.free(cbuf)

	return C.stbi_load_from_memory((*C.uchar)(cbuf), C.int(len(buffer)), &width, &height, &channels, desiredChannels),
		int(width), int(height), channels
}

// ImageFree = stbi_image_free
func ImageFree(image *C.uchar) {
	C.stbi_image_free(unsafe.Pointer(image))
}
