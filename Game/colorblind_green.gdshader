shader_type canvas_item;
render_mode world_vertex_coords;

uniform sampler2D pattern : repeat_enable;
uniform float tile_size;
global uniform vec2 player_position;
varying vec2 world_position;

void vertex() {
	world_position = VERTEX;
}

void fragment() {
	vec2 pattern_pixel_size = 1.0 / vec2(textureSize(pattern, 0));
	vec2 diff = pattern_pixel_size / TEXTURE_PIXEL_SIZE;
	
	vec4 base_texture = texture(TEXTURE, UV);
	vec2 tiled_uv = UV * tile_size * diff;
	vec4 pattern_tex = texture(pattern, tiled_uv * sign(world_position.x - player_position.x) - vec2(TIME*3.0, 0.0));
	COLOR.a = pattern_tex.a * base_texture.a;
}
