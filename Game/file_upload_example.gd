extends Node

@onready var pick_image_file_dialog = $PickImageFileDialog

var web_server_url: String = ProjectSettings.get_setting("application/boardwalk/web_server_url")

# Called when the node enters the scene tree for the first time.
func _ready():
	pick_image_file_dialog.show()
	pick_image_file_dialog.file_selected.connect(_on_file_selected)
	
func _on_file_selected(path: String):
	print("begin upload")
	var req = FileUploader.upload_file(web_server_url + "/api/simplify", path, HTTPClient.METHOD_POST, "image")
	
	req.request_completed.connect(_on_response_received)

func _on_response_received(result: int, response_code: int, headers: PackedStringArray, body: PackedByteArray):
	print("response received %d, %d" % [result, response_code])
	print(headers)
	if response_code != 200:
		var label = Label.new()
		label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
		label.text = body.get_string_from_utf8()
		add_child(label)
		return
	
	for header in headers:
		if header.contains("Content-Type") and header.contains("image/png"):
			var img = Image.new()
			img.load_png_from_buffer(body)
			var tex_rect = TextureRect.new()
			tex_rect.texture = ImageTexture.create_from_image(img)
			add_child(tex_rect)
		else:
			var label = Label.new()
			var json = JSON.parse_string(body.get_string_from_utf8())
			var json_string = JSON.stringify(json, "  ")
			label.text = json_string
			add_child(label)
			return
	
	

func _unhandled_key_input(event):
	if event is InputEventKey and event.is_action_pressed("ui_accept"):
		get_tree().reload_current_scene()
