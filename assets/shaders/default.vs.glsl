#version 410 core

uniform mat4 uProjection;
uniform mat4 uView;
uniform mat4 uModel;

layout(location = 0) in vec3 vertex;

void main() {
    gl_Position = uProjection * uView * uModel * vec4(vertex, 1);
}
