shader_type canvas_item;

uniform sampler2D pattern : repeat_enable;
uniform float tile_size;

void fragment() {
	vec2 pattern_pixel_size = 1.0 / vec2(textureSize(pattern, 0));
	vec2 diff = pattern_pixel_size / TEXTURE_PIXEL_SIZE;
		
	vec4 base_texture = texture(TEXTURE, UV);
	vec2 tiled_uv = UV * tile_size * diff;

	vec4 pattern_tex = texture(pattern, tiled_uv + vec2(0.0, TIME * 3.0));
	COLOR.a = pattern_tex.a * base_texture.a;
}

//void light() {
	// Called for every pixel for every light affecting the CanvasItem.
	// Uncomment to replace the default light processing function with this one.
//}
