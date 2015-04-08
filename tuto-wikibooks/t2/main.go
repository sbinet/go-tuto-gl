package main

import (
	"encoding/binary"
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
)

type context struct {
	w    *glfw.Window
	prog gl.Program
	buf  gl.Buffer
	pos  gl.Attrib
	data []float32
}

func (ctx *context) Delete() {
	gl.DeleteProgram(ctx.prog)
	gl.DeleteBuffer(ctx.buf)
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

	gl.BindBuffer(gl.ARRAY_BUFFER, ctx.buf)
	gl.EnableVertexAttribArray(ctx.pos)
	gl.VertexAttribPointer(ctx.pos, 4, gl.FLOAT, false, 0, 0)
	gl.DrawArrays(gl.TRIANGLES, 0, len(ctx.data))

	// display result
	ctx.w.SwapBuffers()
}

func main() {

	err := glfw.Init()
	if err != nil {
		panic(err)
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

	ctx := context{
		w:   w,
		buf: gl.CreateBuffer(),
		data: []float32{
			+0.0, +0.8, 0, 1,
			-0.8, -0.8, 0, 1,
			+0.8, -0.8, 0, 1,
		},
	}
	ctx.prog, err = glutil.CreateProgram(
		newShader("triangle.v.glsl"),
		newShader("triangle.f.glsl"),
	)
	if err != nil {
		panic(err)
	}
	defer ctx.Delete()

	gl.BindBuffer(gl.ARRAY_BUFFER, ctx.buf)
	gl.BufferData(gl.ARRAY_BUFFER,
		f32.Bytes(binary.LittleEndian, ctx.data...),
		gl.STATIC_DRAW,
	)
	ctx.pos = gl.GetAttribLocation(ctx.prog, "coord")

	ctx.w.SetSizeCallback(onResize)
	ctx.w.SetKeyCallback(onKey)

	gl.UseProgram(ctx.prog)

	for !ctx.w.ShouldClose() {
		display(ctx)
		glfw.PollEvents()
	}

}
