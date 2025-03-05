shader_type canvas_item;

uniform float edge: hint_range(0.0, 1.0) = 0.0;
uniform int steps = 10;

void fragment() {
	int current_step = int(edge * float(steps));
	
	float pos1 = (UV.x + UV.y) / 2.0;
	float pos2 = (UV.x - UV.y + 1.0) / 2.0; 
	COLOR.a = pos1 < float(current_step) / float(steps)
			? 1.0
			: pos1 < float(current_step + 1) / float(steps)
				? (current_step % 2 == 0 ? pos2 : 1.0 - pos2) < (edge - float(current_step) / float(steps)) * float(steps)
					? 1.0
					: 0.0
				: 0.0;
}

//void light() {
//	// Called for every pixel for every light affecting the CanvasItem.
//	// Uncomment to replace the default light processing function with this one.
//}
