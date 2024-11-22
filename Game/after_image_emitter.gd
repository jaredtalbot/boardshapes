@icon("res://icons/afterimageicon.png")
class_name AfterImageEmitter extends Node2D

@export var enabled: bool = true
@export var animated_sprite: AnimatedSprite2D
@export var emit_frequency: float = 0.1
@export var random_colors := PackedColorArray([Color(1, 1, 1, 1)])

var emit_progress: float = 0.0
var current_color: int = 0

func _ready():
	assert(animated_sprite != null)

func _process(delta):
	if enabled:
		emit_progress += delta
		if emit_progress >= emit_frequency:
			emit_progress = fmod(emit_progress, emit_frequency)
			var texture = animated_sprite.sprite_frames.get_frame_texture(animated_sprite.animation, animated_sprite.frame)
			var after_image := AfterImage.new()
			after_image.texture = texture
			after_image.transform = animated_sprite.global_transform
			after_image.flip_h = animated_sprite.flip_h
			after_image.texture_filter = animated_sprite.texture_filter
			after_image.color = random_colors[current_color % len(random_colors)]
			current_color = (current_color + 1) % len(random_colors)
			
			get_tree().current_scene.add_child(after_image)
	else:
		emit_progress = emit_frequency
