package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/sbinet/go-tuto-gl/tuto-wikibooks/glh"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
)

type Point struct {
	x float32
	y float32
}

type context struct {
	w    *glfw.Window
	prog gl.Program

	coord   gl.Attrib
	xoff    gl.Uniform
	xscale  gl.Uniform
	sprite  gl.Uniform
	texture gl.Uniform

	tex gl.Texture
	vbo gl.Buffer

	posbuf  gl.Buffer
	pos     gl.Attrib
	posdata []float32

	colbuf  gl.Buffer
	col     gl.Attrib
	coldata []float32

	eltbuf  gl.Buffer
	elt     gl.Attrib
	eltdata []uint16

	fade  gl.Uniform
	trans gl.Uniform
}

func init() {
	runtime.LockOSThread()
}

func main() {
	w, err := glh.New(width, height, "graph")
	if err != nil {
		panic(err)
	}

	prog, err := glutil.CreateProgram(vtxShader, fragShader)
	if err != nil {
		panic(err)
	}

	ctx := context{
		w:       w,
		prog:    prog,
		coord:   gl.GetAttribLocation(prog, "coord2d"),
		xoff:    gl.GetUniformLocation(prog, "offset_x"),
		xscale:  gl.GetUniformLocation(prog, "scale_x"),
		sprite:  gl.GetUniformLocation(prog, "sprite"),
		texture: gl.GetUniformLocation(prog, "mytexture"),

		tex: gl.CreateTexture(),
		vbo: gl.CreateBuffer(),
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ctx.tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, 15, 15, gl.RGBA, gl.UNSIGNED_BYTE, _res_texture_png,
	)

	gl.BindBuffer(gl.ARRAY_BUFFER, ctx.vbo)

	for !ctx.w.ShouldClose() {
		display(ctx)
		glfw.PollEvents()
	}
}

func display(ctx context) {
}

const (
	width  = 800
	height = 600

	vtxShader = `
#version 120

attribute vec2 coord2d;
varying vec4 f_color;
uniform float offset_x;
uniform float scale_x;
uniform float sprite;

void main(void) {
	gl_Position = vec4((coord2d.x + offset_x) * scale_x, coord2d.y, 0, 1);
	f_color = vec4(coord2d.xy / 2.0 + 0.5, 1, 1);
	gl_PointSize = max(1.0, sprite);
}
`

	fragShader = `
#version 120

uniform sampler2D mytexture;
varying vec4 f_color;
uniform float sprite;

void main(void) {
	if (sprite > 1.0)
		gl_FragColor = texture2D(mytexture, gl_PointCoord) * f_color;
	else
		gl_FragColor = f_color;
}
`
)
