package main

import (
	"fnmatch/fnmatch"
	"strings"
	"unsafe"
)

func main() {}

func decodeString(ptr *uint32, length int) string {
	var str strings.Builder
	pointer := uintptr(unsafe.Pointer(ptr))
	for i := 0; i < length; i++ {
		s := *(*int32)(unsafe.Pointer(pointer + uintptr(i)))
		str.WriteByte(byte(s))
	}

	return str.String()
}

//export _FNM_NOESCAPE
func _FNM_NOESCAPE() int {
	return fnmatch.FNM_NOESCAPE
}

//export _FNM_PATHNAME
func _FNM_PATHNAME() int {
	return fnmatch.FNM_PATHNAME
}

//export _FNM_PERIOD
func _FNM_PERIOD() int {
	return fnmatch.FNM_PERIOD
}

//export _FNM_LEADING_DIR
func _FNM_LEADING_DIR() int {
	return fnmatch.FNM_LEADING_DIR
}

//export _FNM_CASEFOLD
func _FNM_CASEFOLD() int {
	return fnmatch.FNM_CASEFOLD
}

//export _FNM_IGNORECASE
func _FNM_IGNORECASE() int {
	return fnmatch.FNM_IGNORECASE
}

//export _FNM_FILE_NAME
func _FNM_FILE_NAME() int {
	return fnmatch.FNM_FILE_NAME
}

//export alloc
func alloc(size uint32) *byte {
	buf := make([]byte, size)
	return &buf[0]
}

//export Match
func Match(p1 *uint32, l1 int, p2 *uint32, l2 int, flags int) bool {
	a := decodeString(p1, l1)
	b := decodeString(p2, l2)

	return fnmatch.Match(a, b, flags)
}
