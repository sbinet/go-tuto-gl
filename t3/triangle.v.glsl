#version 120

attribute vec4 coord;
attribute vec3 v_color;
varying   vec3 f_color;

void main(void) {
  gl_Position = coord;
  f_color = v_color;
}
