package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
)

const (
	deg2rad = math.Pi / 180.0
)

var (
	start = time.Now()
)

type context struct {
	w    *glfw.Window
	prog gl.Program

	posbuf  gl.Buffer
	pos     gl.Attrib
	posdata []float32

	colbuf  gl.Buffer
	col     gl.Attrib
	coldata []float32

	eltbuf  gl.Buffer
	elt     gl.Attrib
	eltdata []byte

	fade  gl.Uniform
	trans gl.Uniform
}

func (ctx *context) Delete() {
	gl.DeleteProgram(ctx.prog)
	gl.DeleteBuffer(ctx.posbuf)
	gl.DeleteBuffer(ctx.colbuf)
	gl.DeleteBuffer(ctx.eltbuf)
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
	w.SetSizeCallback(onResize)
	w.SetKeyCallback(onKey)

	glfw.SwapInterval(1)

	//	gl.Enable(gl.BLEND)
	//	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.DEPTH_TEST)

	ctx := context{
		w: w,
		posdata: []float32{
			// front
			-0.8, -0.8, +0.8, 1,
			+0.8, -0.8, +0.8, 1,
			+0.8, +0.8, +0.8, 1,
			-0.8, +0.8, +0.8, 1,

			// back
			-0.8, -0.8, -0.8, 1,
			+0.8, -0.8, -0.8, 1,
			+0.8, +0.8, -0.8, 1,
			-0.8, +0.8, -0.8, 1,
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

		eltdata: []byte{
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
	ctx.prog, err = glutil.CreateProgram(
		newShader("cube.v.glsl"),
		newShader("cube.f.glsl"),
	)
	if err != nil {
		panic(err)
	}
	defer ctx.Delete()

	ctx.posbuf = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, ctx.posbuf)
	gl.BufferData(gl.ARRAY_BUFFER,
		f32.Bytes(binary.LittleEndian, ctx.posdata...),
		gl.STATIC_DRAW,
	)
	ctx.pos = gl.GetAttribLocation(ctx.prog, "coord")

	ctx.colbuf = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, ctx.colbuf)
	gl.BufferData(gl.ARRAY_BUFFER,
		f32.Bytes(binary.LittleEndian, ctx.coldata...),
		gl.STATIC_DRAW,
	)

	ctx.eltbuf = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ctx.eltbuf)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER,
		ctx.eltdata,
		gl.STATIC_DRAW,
	)

	ctx.col = gl.GetAttribLocation(ctx.prog, "v_color")
	//ctx.fade = gl.GetUniformLocation(ctx.prog, "fade")
	ctx.trans = gl.GetUniformLocation(ctx.prog, "m_transform")

	for !ctx.w.ShouldClose() {
		display(ctx)
		glfw.PollEvents()

	}

}

func display(ctx context) {
	// clear the background as black
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(ctx.prog)

	// // 1<->+1 every 5 seconds
	// tx := f32.Sin(float32(time.Since(start).Seconds()) * (2 * float32(math.Pi)) / 5.0)
	//
	// // 45-degrees per second
	// angle := float32(time.Since(start).Seconds()) * 45 * deg2rad
	//
	// m := f32.Mat4{
	//         {1, 1, 1, 1},
	//         {1, 1, 1, 1},
	//         {1, 1, 1, 1},
	//         {1, 1, 1, 1},
	// }
	// m.Translate(&m, tx, 0, 0)
	// m.Rotate(&m, f32.Radian(angle), &f32.Vec3{0, 0, 1})
	//
	// gl.UniformMatrix4fv(ctx.trans, flatten(&m))

	gl.EnableVertexAttribArray(ctx.col)
	gl.EnableVertexAttribArray(ctx.pos)
	gl.EnableVertexAttribArray(ctx.elt)

	gl.BindBuffer(gl.ARRAY_BUFFER, ctx.colbuf)
	gl.VertexAttribPointer(ctx.col, 3, gl.FLOAT, false, 0, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, ctx.posbuf)
	gl.VertexAttribPointer(ctx.pos, 4, gl.FLOAT, false, 0, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ctx.eltbuf)
	sz := gl.GetBufferParameteri(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE)
	gl.DrawElements(gl.TRIANGLE_STRIP, sz, gl.FLOAT, 0)

	gl.DisableVertexAttribArray(ctx.col)
	gl.DisableVertexAttribArray(ctx.pos)
	gl.DisableVertexAttribArray(ctx.elt)

	// display result
	ctx.w.SwapBuffers()
}

func flatten(m *f32.Mat4) []float32 {
	o := make([]float32, 0, 16)
	o = append(o, (*m)[0][:]...)
	o = append(o, (*m)[1][:]...)
	o = append(o, (*m)[2][:]...)
	o = append(o, (*m)[3][:]...)
	return o
}

func init() {
	runtime.LockOSThread()
}
