extends Node

func _ready():
	AccessibilityShaderManager.apply_shaders()

func _on_texture_button_pressed():
	ScreenTransitioner.change_scene_to_file("res://menus/main.tscn")
