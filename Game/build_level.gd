extends Node

@onready var pick_image_file_dialog = $PickImageFileDialog

var web_server_url: String = ProjectSettings.get_setting("application/boardwalk/web_server_url")

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	pick_image_file_dialog.show()
	pick_image_file_dialog.file_selected.connect(_on_file_selected)

func _on_file_selected(path: String):
	print("begin upload")
	var req = FileUploader.upload_file(web_server_url + "/api/build-level", path, HTTPClient.METHOD_POST, "image")
	
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
			return
		else:		
			var json = JSON.parse_string(body.get_string_from_utf8())
			var json_string = JSON.stringify(json, "  ")
			generate_nodes(json_string)
			return
	

func generate_nodes(json_string: String):
	var json = JSON.parse_string(json_string)
	var level = Node.new()
	for item in json:
		var region = Node2D.new()
		var byte_pool = Marshalls.base64_to_raw(item["regionImage"])
		var img = Image.new()
		img.load_png_from_buffer(byte_pool)
		var tex_rect = TextureRect.new()
		tex_rect.texture = ImageTexture.create_from_image(img)
		region.add_child(tex_rect)
		var collision = CollisionPolygon2D.new()
		var mesh = item["mesh"] as Array
		var vectormesh = mesh.map(func(v: Dictionary): return Vector2(v["x"], v["y"]))
		collision.polygon = vectormesh
		var col = StaticBody2D.new()
		col.add_child(collision)
		region.add_child(col)
		region.position = Vector2(item["cornerX"], item["cornerY"])
		level.add_child(region)
	add_child(level)
