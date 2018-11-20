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
    vec3 ambient = texture(uAmbientMap, p_TexCoord).rgb;
    vec3 diffuse = texture(uDiffuseMap, p_TexCoord).rgb;
    vec3 specular = texture(uSpecularMap, p_TexCoord).rgb;

    vec3 normal = normalize(p_Normal.xyz);
    diffuse *= max(dot(normal, p_LightDir), 0.0);

    vec3 halfway = normalize(p_LightDir + p_ViewDir);
    specular *= pow(max(dot(normal, halfway), 0.0), 16.0) * 0.5;

    _Color = vec4(ambient + diffuse + specular, 1.0);
}
