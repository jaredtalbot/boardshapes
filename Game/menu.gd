extends Control

func _ready():
	if OS.has_feature("web"):
		$MarginContainer/VBoxContainer/Exit.hide()

func _on_start_pressed():
	get_tree().change_scene_to_file("res://start_menu.tscn")


func _on_exit_pressed():
	get_tree().quit()
