uniform sampler2D uTexture; 

in vec2 p_TexCoord;

out vec4 _Color;

void main() {
    _Color = texture(uTexture, p_TexCoord);
}