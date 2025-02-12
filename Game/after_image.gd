@icon("res://icons/afterimageicon.png")
class_name AfterImage extends Sprite2D

var color: Color

func _ready():
	z_index = -1
	material = ShaderMaterial.new()
	material.shader = preload("res://shaders/silhouette.gdshader")
	material.set_shader_parameter("color", color)
	var tween := create_tween()
	modulate = Color(1, 1, 1, 0.5)
	tween.tween_property(self, "modulate", Color(1, 1, 1, 0), 0.75)
	tween.tween_callback(queue_free)
