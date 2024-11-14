@icon("res://icons/afterimageicon.png")
class_name AfterImage extends Sprite2D

func _ready():
	z_index = -1
	var tween := create_tween()
	var new_color = Color(modulate)
	new_color.a = 0
	tween.tween_property(self, "modulate", new_color, 0.5)
	tween.tween_callback(queue_free)
