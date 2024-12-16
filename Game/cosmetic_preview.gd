class_name CosmeticPreview extends Node2D

static var animation_list := PackedStringArray([
	&"idle animation",
	&"running",
	&"jumping",
	&"sliding",
	&"dash",
])

@export var animation_length := 2.0
var next_animation_progress := 0.0
var current_animation_index := 0

@onready var animation_player = $AnimationPlayer
@onready var hat_pos = $HatPivot/HatPos

func _ready():
	animation_player.play(animation_list[current_animation_index])

func _process(delta):
	next_animation_progress += delta
	if next_animation_progress > animation_length:
		current_animation_index = (current_animation_index + 1) % len(animation_list)
		animation_player.play(animation_list[current_animation_index])
		next_animation_progress = 0.0

func equip_hat(hat: PackedScene):
	assert(hat_pos.get_child_count() < 2)
	if hat != null:
		if hat_pos.get_child_count() > 0:
			var existing_hat := hat_pos.get_child(0)
			existing_hat.replace_by(hat.instantiate())
			existing_hat.queue_free()
		else:
			hat_pos.add_child(hat.instantiate())
	else:
		if hat_pos.get_child_count() > 0:
			hat_pos.get_child(0).queue_free()
