shader_type canvas_item;

void fragment() {
	vec4 c = texture(TEXTURE, UV);
	COLOR = vec4(1.0 - c.rgb, c.a);
}
