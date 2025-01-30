extends Control

func _ready():
	Music.stop_all_layers()
	$CreditsScreen.hide()
	$Back.hide()
	if OS.has_feature("web"):
		$MarginContainer/VBoxContainer/Exit.hide()
	if RenderingServer.get_default_clear_color() == Color(0, 0, 0, 1):
		$OptionsWindow/VBoxContainer/DarkMode.set_pressed_no_signal(true)
	$OptionsWindow/VBoxContainer/ColorblindMode.set_pressed_no_signal(ProjectSettings.get_setting("rendering/environment/defaults/color_blind_mode"))
	$MarginContainer/VBoxContainer/Start.call_deferred("grab_focus")
	
func _on_start_pressed():
	get_tree().change_scene_to_file("res://menus/start_menu.tscn")


func _on_exit_pressed():
	get_tree().quit()

func _on_options_pressed():
	$OptionsWindow.show()
	
func _on_options_close_pressed():
	$OptionsWindow.hide()

func _on_credits_pressed():
	$CreditsScreen.show()
	$Back.show()


func _on_back_pressed() -> void:
	$Back.hide()
	$CreditsScreen.hide()
