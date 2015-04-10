package main

import (
	"io/ioutil"
)

func newShader(fname string) string {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	return string(buf)
}
