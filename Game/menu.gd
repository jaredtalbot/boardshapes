extends Control

func _ready():
	if OS.has_feature("web"):
		$MarginContainer/VBoxContainer/Exit.hide()
	if RenderingServer.get_default_clear_color() == Color(0, 0, 0, 1):
		$OptionsWindow/DarkMode.set_pressed_no_signal(true)
	
func _on_start_pressed():
	get_tree().change_scene_to_file("res://start_menu.tscn")


func _on_exit_pressed():
	get_tree().quit()

func _on_options_pressed():
	$OptionsWindow.show()

func _on_options_close_pressed():
	$OptionsWindow.hide()
	
func _on_dark_mode_toggled(toggled):
	if toggled:
		RenderingServer.set_default_clear_color(Color(0, 0, 0, 1))
	else:
		RenderingServer.set_default_clear_color(Color(1, 1, 1, 1))
