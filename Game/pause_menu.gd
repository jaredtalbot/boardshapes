extends Window

signal resume_button_pressed
signal quit_button_pressed

func _input(event):
	if event.is_action_pressed("pause"):
		resume_button_pressed.emit()
		get_viewport().set_input_as_handled()

func _on_resume_button_pressed():
	resume_button_pressed.emit()

func _on_quit_button_pressed():
	quit_button_pressed.emit()
