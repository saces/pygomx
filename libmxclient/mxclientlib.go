package main

import "C"

//export hello
func hello(name *C.char) *C.char {
	goName := C.GoString(name)
	result := "Hello " + goName
	return C.CString(result)
}

func main() {}
