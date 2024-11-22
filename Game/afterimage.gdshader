shader_type canvas_item;

uniform vec3 color;

void fragment() {
	COLOR.rgba = vec4(color, COLOR.a);
}
