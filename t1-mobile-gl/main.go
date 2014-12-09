package main

import (
	"fmt"

	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
	"golang.org/x/mobile/gl"
)

const (
	vshader = `
#version 120
attribute vec2 coord2d;
void main(void) {
  gl_Position = vec4(coord2d, 0.0, 1.0);
}
`
	fshader = `
#version 120
void main(void) {
  gl_FragColor[0] = 0.0;
  gl_FragColor[1] = 0.0;
  gl_FragColor[2] = 1.0;
}`
)

func onError(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int,
	action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

func onResize(w *glfw.Window, width, height int) {
	gl.Viewport(0, 0, width, height)

	//heightUnif.Uniform1i(height)
}

func initOpenGl(window *glfw.Window, w, h int) {
	w, h = window.GetSize() // query window to get screen pixels
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, width, height)
	gl.ClearColor(.25, .88, .83, 1) // turquoise
}

func main() {
	glfw.SetErrorCallback(onError)

	if !glfw.Init() {
		panic("init glfw")
	}
	defer glfw.Terminate()

	w, err := glfw.CreateWindow(640, 480, "my first triangle", nil, nil)
	if err != nil {
		panic(err)
	}
	defer w.Destroy()
	w.SetKeyCallback(onKey)

	w.MakeContextCurrent()
	glfw.SwapInterval(1)

	//gl.Init()

	prog := glh.NewProgram(
		glh.Shader{gl.VERTEX_SHADER, vshader},
		glh.Shader{gl.FRAGMENT_SHADER, fshader},
	)
	defer prog.Delete()

	prog.Link()

	prog.Use()
	w.SetSizeCallback(onResize)
	w.GetSize()

	for !w.ShouldClose() {
		w.SwapBuffers()
		glfw.PollEvents()
	}

	//gl.ProgramUnuse()
}
