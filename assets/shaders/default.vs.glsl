#version 300 es

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

layout(location = 0) in vec3 vertex;

void main() {
    gl_Position = projection * camera * model * vec4(vertex, 1);
}
