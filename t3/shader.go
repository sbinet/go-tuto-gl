package main

import (
	"io/ioutil"

	"github.com/go-gl/gl"
	"github.com/go-gl/glh"
)

func NewShader(typ gl.GLenum, fname string) (glh.Shader, error) {
	var shader glh.Shader
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return shader, err
	}

	return glh.Shader{typ, string(buf)}, nil
}

func MustShader(typ gl.GLenum, fname string) glh.Shader {
	shader, err := NewShader(typ, fname)
	if err != nil {
		panic(err)
	}
	return shader
}
