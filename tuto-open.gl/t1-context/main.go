package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
)

func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	win, err := glfw.CreateWindow(800, 600, "OpenGL", nil, nil)
	if err != nil {
		panic(err)
	}

	win.MakeContextCurrent()

	for !win.ShouldClose() {
		win.SwapBuffers()
		glfw.PollEvents()

		if win.GetKey(glfw.KeyEscape) == glfw.Press {
			win.SetShouldClose(true)
		}
	}

	glfw.Terminate()
}
