uniform sampler2D uDiffuseMap; 

in vec2 p_TexCoord;

out vec4 _Color;

void main() {
    _Color = texture(uDiffuseMap, p_TexCoord);
}