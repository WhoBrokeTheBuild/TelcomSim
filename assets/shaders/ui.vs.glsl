uniform mat4 uProjection;

layout(location = 0) in vec3 _Position;
layout(location = 2) in vec2 _TexCoord;

out vec2 p_TexCoord;

void main() {
    p_TexCoord = vec2(_TexCoord.x, 1.0 - _TexCoord.y);

    gl_Position = uProjection * vec4(_Position, 1);
}
