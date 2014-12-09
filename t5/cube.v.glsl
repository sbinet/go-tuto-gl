#version 120

attribute vec4 coord;
attribute vec3 v_color;
varying   vec3 f_color;
uniform   mat4 m_transform;

void main(void) {
  gl_Position = m_transform * coord;
  f_color = v_color;
}
