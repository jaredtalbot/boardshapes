shader_type canvas_item;

uniform sampler2D pattern : repeat_enable;
uniform float tile_size;

void fragment() {
	vec2 pattern_pixel_size = 1.0 / vec2(textureSize(pattern, 0));
	vec2 diff = pattern_pixel_size / TEXTURE_PIXEL_SIZE;

	vec4 in_tex = texture(TEXTURE, UV);

	vec2 centered_uv = UV * 2.0 - 1.0;
	vec2 tiled_uv = centered_uv * tile_size * diff;

	tiled_uv  = tiled_uv * 0.5 + 0.5;

	vec4 pattern_tex = texture(pattern, tiled_uv);
	COLOR.a = pattern_tex.a * in_tex.a;
}