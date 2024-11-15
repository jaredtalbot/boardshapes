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
		var error_dialog = ErrorDialog.new()
		error_dialog.set_text_to_error_message(body, response_code)
		error_dialog.show()
		add_child(error_dialog)
		error_dialog.confirmed.connect(go_back)
		error_dialog.canceled.connect(go_back)
		return
	
	var level_data = body.get_string_from_utf8()
	var generated_level = level_generator.generate_nodes(level_data)
	if generated_level == null:
		var error_dialog = ErrorDialog.new()
		error_dialog.set_text_to_error_message("Could not generate level with server response.")
		error_dialog.show()
		add_child(error_dialog)
		error_dialog.confirmed.connect(go_back)
		error_dialog.canceled.connect(go_back)
		return
		
	add_child(generated_level)
	add_player()
	loading_indicator.hide()
	get_tree().paused = true
	$StartEndSelection/StartSelect.disabled = false
	$StartEndSelection/StartSelect.show()

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
			$AudioStreamPlayer.play()
			ray_cast.queue_free()
			return
	ray_cast.queue_free()
	
func _on_quit_button_pressed():
	get_node("./QuitMenu/QuitWindow").show()
	get_tree().paused = true

func _on_no_button_pressed():
	get_node("./QuitMenu/QuitWindow").hide()
	get_tree().paused = false

func _on_back_button_pressed():
	get_tree().paused = false
	go_back()
	
func go_back():
	get_tree().change_scene_to_file("res://start_menu.tscn")

func _on_exit_to_main_menu_button_pressed():
	get_tree().change_scene_to_file("res://main.tscn")

func _set_player_start():
	var player = $Player
	player.initial_position = get_viewport().get_mouse_position()
	player.position = player.initial_position
	$StartEndSelection/StartSelect.disabled = true
	$StartEndSelection/StartSelect.hide()
	$StartEndSelection/EndSelect.disabled = false
	$StartEndSelection/EndSelect.show()

func _set_goal_position():
	var goal = $Goal
	goal.position = get_viewport().get_mouse_position()
	$StartEndSelection/EndSelect.disabled = true
	$StartEndSelection/EndSelect.hide()
	$TouchScreenControls.show()
	$Goal.show()
	get_tree().paused = false
	
func _goal_reached(player: Node2D):
	player.set_physics_process(false)
	$VictoryScreen/Victory.show()
	
func _on_audio_stream_player_finished():
	$AudioStreamPlayer.play()

func _on_restart_button_pressed():
	$VictoryScreen/Victory.hide()
	$Goal.hide()
	
	loading_indicator.hide()
	get_tree().paused = true
	$StartEndSelection/StartSelect.disabled = false
	$StartEndSelection/StartSelect.show()
	var player = $Player
	player.set_physics_process(true)
