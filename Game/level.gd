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
	add_player()
	loading_indicator.hide()

func add_player():
	var ray_cast = RayCast2D.new()
	ray_cast.position = Vector2(0, -50)
	ray_cast.target_position = Vector2(0, 1080)
	add_child(ray_cast)
	for i in range(1920):
		ray_cast.position.x = i
		ray_cast.force_raycast_update()
		if ray_cast.is_colliding():
			var player = preload("res://player.tscn").instantiate()
			player.position = Vector2(i + player.get_node("CollisionShape2D").shape.get_rect().size.x/2, -100)
			add_child(player)
			ray_cast.queue_free()
			return
	ray_cast.queue_free()
	
	
