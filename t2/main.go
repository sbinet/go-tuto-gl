package main

import (
	"fmt"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
)

type context struct {
	w    *glfw.Window
	prog gl.Program
	pos  gl.Buffer
	data []float32
}

func (ctx *context) Delete() {
	ctx.prog.Delete()
	ctx.pos.Delete()
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

func display(ctx context) {
	// clear the background as black
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	ctx.pos.Bind(gl.ARRAY_BUFFER)
	coord := gl.AttribLocation(ctx.prog.GetAttribLocation("coord"))
	coord.EnableArray()

	coord.AttribPointer(
		4,        // number of elements per vertex: (x,y)
		gl.FLOAT, // type of each element
		false,    // take our values as-is
		0,
		uintptr(0),
	)
	const sz = 4 // size of float32 in bytes
	gl.DrawArrays(gl.TRIANGLES, 0, len(ctx.data)/sz)
	coord.DisableArray()

	// display result
	ctx.w.SwapBuffers()
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

	w.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Init()

	ctx := context{
		w: w,
		data: []float32{
			+0.0, +0.8, 0, 1,
			-0.8, -0.8, 0, 1,
			+0.8, -0.8, 0, 1,
		},
	}
	ctx.pos = genVertexBuffer(ctx.data)
	ctx.prog = glh.NewProgram(
		MustShader(gl.VERTEX_SHADER, "triangle.v.glsl"),
		MustShader(gl.FRAGMENT_SHADER, "triangle.f.glsl"),
	)
	defer ctx.Delete()

	ctx.prog.Use()
	ctx.w.SetSizeCallback(onResize)
	ctx.w.SetKeyCallback(onKey)

	ctx.prog.Link()

	for !ctx.w.ShouldClose() {
		display(ctx)
		glfw.PollEvents()

	}

	gl.ProgramUnuse()
}

func genVertexBuffer(vtx []float32) gl.Buffer {
	const sz = 4 // size of float32 in bytes

	buffer := gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(vtx)*sz, vtx, gl.STATIC_DRAW)

	buffer.Unbind(gl.ARRAY_BUFFER)
	return buffer
}
