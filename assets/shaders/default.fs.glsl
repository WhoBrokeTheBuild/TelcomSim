uniform vec3 uAmbient;
uniform vec3 uDiffuse;
uniform vec3 uSpecular;

uniform sampler2D uAmbientMap; 
uniform sampler2D uDiffuseMap; 
uniform sampler2D uSpecularMap; 

in vec4 p_Position;
in vec4 p_Normal;
in vec2 p_TexCoord;

in vec3 p_LightDir;
in vec3 p_ViewDir;

out vec4 _Color;

void main() {
    vec3 diffuse = texture(uDiffuseMap, p_TexCoord).rgb;
    vec3 ambient = uAmbient;
    vec3 specular = uSpecular;

    vec4 normal = normalize(p_Normal);

    float diff = max(dot(normal.xyz, p_LightDir), 0.0);
    diffuse *= diff;

    _Color = vec4(texture(uDiffuseMap, p_TexCoord).rgb, 1.0);
}
