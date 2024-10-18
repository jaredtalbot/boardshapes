class_name ErrorDialog extends AcceptDialog

func _init():
	title = "Error"
	size = Vector2(500, 400)
	initial_position = WINDOW_INITIAL_POSITION_CENTER_PRIMARY_SCREEN
	popup_window = true

func set_text_to_error_message(body: Variant, error_code: int = 0):
	if body is PackedByteArray:
		body = body.get_string_from_utf8()
	if body is String:
		body = JSON.parse_string(body)
	
	if body is Dictionary and "errorMessage" in body:
		dialog_text = body["errorMessage"]
	else:
		dialog_text = "An unknown error has occured."
	
	if error_code != 0:
		dialog_text += "\nError Code: %d" % error_code
