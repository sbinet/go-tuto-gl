package main

import (
	"fmt"

	"github.com/go-gl-legacy/glh"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type context struct {
	w       *glfw.Window
	prog    gl.Program
	pos     gl.Buffer
	col     gl.Buffer
	elt     gl.Buffer
	fade    gl.UniformLocation
	trans   gl.UniformLocation
	posdata []float32
	coldata []float32
	eltdata []int8
}

func (ctx *context) Delete() {
	ctx.prog.Delete()
	ctx.pos.Delete()
	ctx.col.Delete()
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

	// ctx.fade = ctx.prog.GetUniformLocation("fade")
	// ctx.fade.Uniform1f(0.5)

	ctx.trans = ctx.prog.GetUniformLocation("m_transform")
	ctx.trans.UniformMatrix4fv(false, [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})

	ctx.col.Bind(gl.ELEMENT_ARRAY_BUFFER)
	colors := gl.AttribLocation(ctx.prog.GetAttribLocation("v_color"))

	ctx.pos.Bind(gl.ELEMENT_ARRAY_BUFFER)
	coord := gl.AttribLocation(ctx.prog.GetAttribLocation("coord"))

	coord.EnableArray()
	colors.EnableArray()

	coord.AttribPointer(
		4,          // number of elements per vertex: (x,y)
		gl.FLOAT,   // type of each element
		false,      // take our values as-is
		0,          // no extra data between each position
		uintptr(0), // offset of first element
	)

	colors.AttribPointer(
		3,          // number of elements per vertex, here (r,g,b)
		gl.FLOAT,   // the type of each element
		false,      // take our values as-is
		0,          // no extra data between each position
		uintptr(0), // offset of first element
	)

	sz := int(gl.GetBufferParameteriv(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE))
	ctx.elt.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.DrawElements(gl.TRIANGLES, sz, gl.UNSIGNED_SHORT, ctx.eltdata)

	coord.DisableArray()
	colors.DisableArray()

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
		posdata: []float32{
			// front
			-1.0, -1.0, 1.0, 1,
			+1.0, -1.0, 1.0, 1,
			+1.0, +1.0, 1.0, 1,
			-1.0, +1.0, 1.0, 1,
			// back
			-1.0, -1.0, -1.0, 1,
			+1.0, -1.0, -1.0, 1,
			+1.0, +1.0, -1.0, 1,
			-1.0, +1.0, -1.0, 1,
		},
		coldata: []float32{
			// front colors
			1.0, 0.0, 0.0,
			0.0, 1.0, 0.0,
			0.0, 0.0, 1.0,
			1.0, 1.0, 1.0,
			// back colors
			1.0, 0.0, 0.0,
			0.0, 1.0, 0.0,
			0.0, 0.0, 1.0,
			1.0, 1.0, 1.0,
		},
		eltdata: []int8{
			// front
			0, 1, 2,
			2, 3, 0,
			// top
			3, 2, 6,
			6, 7, 3,
			// back
			7, 6, 5,
			5, 4, 7,
			// bottom
			4, 5, 1,
			1, 0, 4,
			// left
			4, 0, 3,
			3, 7, 4,
			// right
			1, 5, 6,
			6, 2, 1,
		},
	}
	ctx.pos = genFloatBuffer(ctx.posdata)
	ctx.col = genFloatBuffer(ctx.coldata)
	ctx.elt = genIntBuffer(ctx.eltdata)
	ctx.prog = glh.NewProgram(
		MustShader(gl.VERTEX_SHADER, "cube.v.glsl"),
		MustShader(gl.FRAGMENT_SHADER, "cube.f.glsl"),
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

func genFloatBuffer(data []float32) gl.Buffer {
	const sz = 4 // size of float32 in bytes

	buffer := gl.GenBuffer()
	buffer.Bind(gl.ELEMENT_ARRAY_BUFFER)

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data)*sz, data, gl.STATIC_DRAW)

	buffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)
	return buffer
}

func genIntBuffer(data []int8) gl.Buffer {

	buffer := gl.GenBuffer()
	buffer.Bind(gl.ELEMENT_ARRAY_BUFFER)

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data), data, gl.STATIC_DRAW)

	buffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)
	return buffer
}
