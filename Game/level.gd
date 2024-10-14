extends Node

var base_url = ProjectSettings.get_setting("application/boardwalk/web_server_url")

@onready var level_generator = $LevelGenerator
@onready var loading_indicator = $LoadingIndicator

func create_level(img: Image):
	loading_indicator.show()
	loading_indicator.set_text("Uploading Level...")
	var buffer = img.save_png_to_buffer()
	var request = FileUploader.upload_buffer(base_url + "/api/build-level", buffer, "image.png", HTTPClient.METHOD_POST, "image")
	request.request_completed.connect(_on_response_received)

func _on_response_received(result: int, response_code: int, headers: PackedStringArray, body: PackedByteArray):
	if response_code != HTTPClient.RESPONSE_OK:
		#todo: error handling!!!
		return
	
	var level_data = body.get_string_from_utf8()
	var generated_level = level_generator.generate_nodes(level_data)
	if generated_level == null:
		JavaScriptBridge.eval("console.log(\"level generation was not successful\")", true)
		#todo: error handling!!!
		return
		
	add_child(generated_level)
	loading_indicator.hide()
