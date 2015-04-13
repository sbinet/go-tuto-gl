package glh

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
	"golang.org/x/mobile/gl"
)

func New(width, height int, title string) (*glfw.Window, error) {
	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	w, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	w.MakeContextCurrent()
	w.SetSizeCallback(onResize)
	w.SetKeyCallback(onKey)

	glfw.SwapInterval(1)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.DEPTH_TEST)

	return w, err
}

func onError(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int,
	action glfw.Action, mods glfw.ModifierKey) {
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	}
}

func onResize(window *glfw.Window, w, h int) {
	gl.Viewport(0, 0, w, h)
}
