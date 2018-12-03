uniform vec4 uAmbient;
uniform vec4 uDiffuse;
uniform vec4 uSpecular;

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
    vec4 ambient = uAmbient + texture(uAmbientMap, p_TexCoord);
    vec4 diffuse = uDiffuse + texture(uDiffuseMap, p_TexCoord);
    vec4 specular = uSpecular + texture(uSpecularMap, p_TexCoord);

    vec3 normal = normalize(p_Normal.xyz);
    diffuse *= max(dot(normal, p_LightDir), 0.0);

    vec3 halfway = normalize(p_LightDir + p_ViewDir);
    specular *= pow(max(dot(normal, halfway), 0.0), 16.0) * 0.5;

    _Color = texture(uDiffuseMap, p_TexCoord);
}
