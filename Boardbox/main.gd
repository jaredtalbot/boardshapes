extends Node

var base_url = ProjectSettings.get_setting("application/boardbox/web_server_url")

@onready var pick_image_file_dialog = $PickImageFileDialog
@onready var web_pick_image_file = $WebPickImageFile
@onready var loading_indicator = $Other/LoadingIndicator
@onready var error_dialog = $ErrorDialog
@onready var shape_generator = $ShapeGenerator
@onready var shapes = $Shapes


func _on_upload_image_button_pressed():
	if OS.has_feature("web"):
		web_pick_image_file.show()
	else:
		pick_image_file_dialog.show()

func _on_pick_image_file_dialog_file_selected(path):
	var img = Image.load_from_file(path)
	if img == null:
		return
	
	create_level(img)
	

func _on_web_pick_image_file_file_loaded(content: PackedByteArray, filename: String):
	var image = Image.new()
	var error: Error
	if filename.ends_with("png"):
		error = image.load_png_from_buffer(content)
	elif filename.ends_with("jpg") or filename.ends_with("jpeg"):
		error = image.load_jpg_from_buffer(content)
	
	if error != 0:
		return
	
	create_level(image)

func create_level(img: Image):
	loading_indicator.show()
	loading_indicator.set_text("Uploading Image...")
	var buffer = img.save_png_to_buffer()
	var request = FileUploader.upload_buffer(base_url + "/api/build-level", buffer, "image.png", HTTPClient.METHOD_POST, "image")
	request.request_completed.connect(_on_response_received)

func _on_response_received(result: int, response_code: int, headers: PackedStringArray, body: PackedByteArray):
	if response_code != HTTPClient.RESPONSE_OK:
		error_dialog.set_text_to_error_message(body, response_code)
		error_dialog.show()
		loading_indicator.hide()
		return
	
	var level_data = body.get_string_from_utf8()
	var generated_level = shape_generator.generate_nodes(level_data)
	if generated_level == null:
		error_dialog.set_text_to_error_message("Could not generate level with server response.")
		error_dialog.show()
		loading_indicator.hide()
		return
	
	shapes.add_child(generated_level)
	loading_indicator.hide()
