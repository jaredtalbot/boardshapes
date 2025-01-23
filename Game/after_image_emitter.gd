@icon("res://icons/afterimageicon.png")
class_name AfterImageEmitter extends Node2D

@export var sprite: Node2D
@export_group("Auto Emit")
@export var auto_emit: bool = false
@export var emit_frequency: float = 0.1
@export var random_colors := PackedColorArray([Color(1, 1, 1, 1)])

var emit_progress: float = 0.0
var current_color: int = 0

signal emitted_after_image(color: Color)

func _ready():
	assert(sprite != null, "Make sure to set the sprite node to create after images for")

func _process(delta):
	if auto_emit:
		emit_progress += delta
		if emit_progress >= emit_frequency:
			emit_progress = fmod(emit_progress, emit_frequency)
			emit_after_image(random_colors[current_color % len(random_colors)])
			current_color = (current_color + 1) % len(random_colors)
	else:
		emit_progress = emit_frequency

func emit_after_image(color: Color = Color.WHITE):
	var texture : Texture2D
	if sprite is Sprite2D:
		texture = sprite.texture
	elif sprite is AnimatedSprite2D:
		texture = sprite.sprite_frames.get_frame_texture(sprite.animation, sprite.frame)
	var after_image := AfterImage.new()
	after_image.texture = texture
	after_image.transform = sprite.global_transform
	after_image.flip_h = sprite.flip_h
	after_image.texture_filter = sprite.texture_filter
	after_image.color = color
	get_tree().current_scene.add_child(after_image)
	emitted_after_image.emit(color)
