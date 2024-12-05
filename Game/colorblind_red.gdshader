shader_type canvas_item;

uniform sampler2D pattern : repeat_enable;
uniform float tile_size;

void fragment() {
	vec2 pattern_pixel_size = 1.0 / vec2(textureSize(pattern, 0));
	vec2 diff = pattern_pixel_size / TEXTURE_PIXEL_SIZE;
		
	vec4 base_texture = texture(TEXTURE, UV);
	vec2 tiled_uv = UV * tile_size * diff;

	vec4 pattern_tex = texture(pattern, tiled_uv);
	COLOR.a = (1.0 - (1.0 - pattern_tex.a) * (sin(TIME * TAU * 2.0)/4.0 + 3.0/4.0)) * base_texture.a;
}