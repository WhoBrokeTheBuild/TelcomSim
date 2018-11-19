uniform mat4 uProjection;
uniform mat4 uView;
uniform mat4 uModel;

uniform vec3 uLight;
uniform vec3 uCamera;

layout(location = 0) in vec3 _Position;
layout(location = 1) in vec3 _Normal;
layout(location = 2) in vec2 _TexCoord;

out vec4 p_Position;
out vec4 p_Normal;
out vec2 p_TexCoord;

out vec3 p_LightDir;
out vec3 p_ViewDir;

void main() {
    p_Position = uModel * vec4(_Position, 1.0);
    p_Normal   = uModel * vec4(_Normal, 1.0);
    p_TexCoord = vec2(_TexCoord.x, 1.0 - _TexCoord.y);

    p_LightDir = normalize(uLight - p_Position.xyz);
    p_ViewDir  = normalize(uCamera - p_Position.xyz);

    gl_Position = uProjection * uView * uModel * vec4(_Position, 1);
}
