extends Node

func _on_dark_mode_toggled(toggled):
	if toggled:
		RenderingServer.set_default_clear_color(Color(0, 0, 0, 1))
	else:
		RenderingServer.set_default_clear_color(Color(1, 1, 1, 1))

func _on_colorblind_mode_toggled(toggled):
	if toggled:
		ProjectSettings.set_setting("rendering/environment/defaults/color_blind_mode", true)
	else:
		ProjectSettings.set_setting("rendering/environment/defaults/color_blind_mode", false)
