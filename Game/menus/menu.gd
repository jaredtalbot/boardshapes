extends Control

func _ready():
	for hat in Unlocks.HAT_LIST.data:
		var hat_path = hat.get("path")
		if hat_path != null:
			ResourceLoader.load(hat_path)
		
	Music.stop_all_layers()
	$MarginContainer/VBoxContainer/Start.call_deferred("grab_focus")
	
func _on_start_pressed():
	ScreenTransitioner.change_scene_to_file("res://menus/start_menu.tscn")


func _on_exit_pressed():
	if OS.has_feature("web"):
		JavaScriptBridge.eval("window.close()");
	else:
		get_tree().quit()

func _on_options_pressed():
	$OptionsWindow.show()
	
func _on_options_close_pressed():
	$OptionsWindow.hide()

func _on_credits_pressed():
	ScreenTransitioner.change_scene_to_file("res://credits_screen.tscn")


func _on_back_pressed() -> void:
	$Back.hide()
	$CreditsScreen.hide()
