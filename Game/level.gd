extends Node

var base_url = ProjectSettings.get_setting("application/boardwalk/web_server_url")

@onready var level_generator = $LevelGenerator
@onready var loading_indicator = $LoadingIndicator
@onready var multiplayer_timer = $MultiplayerTimer
@onready var multiplayer_controller = $MultiplayerController

func _ready():
	$"QuitMenu/QuitWindow/volume-slider".set_value_no_signal(100)

func create_level(img: Image, options: Dictionary):
	loading_indicator.show()
	loading_indicator.set_text("Uploading Level...")
	if RenderingServer.get_default_clear_color() == Color(0, 0, 0, 1):
		loading_indicator.set_text_color(Color(1, 1, 1, 1))
	var buffer = img.save_png_to_buffer()
	var request = FileUploader.upload_buffer(base_url + "/api/build-level", buffer, "image.png", HTTPClient.METHOD_POST, "image", options)
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
	if RenderingServer.get_default_clear_color() == Color(0, 0, 0, 1):
		$QuitButton.material = ShaderMaterial.new()
		$QuitButton.material.shader = load("res://color_invert.gdshader")
	add_player()
	multiplayer_controller.try_connect(str(hash(level_data)))
	loading_indicator.hide()
	get_tree().paused = true
	$StartEndSelection/StartSelect.disabled = false
	$StartEndSelection/StartSelect.show()

func add_player():
	var player = preload("res://player.tscn").instantiate()
	if RenderingServer.get_default_clear_color() == Color(0, 0, 0, 1):
		player.get_node("AnimatedSprite2D").material = ShaderMaterial.new()
		player.get_node("AnimatedSprite2D").material.shader = load("res://color_invert.gdshader")
	add_child(player)
	
func _on_quit_button_pressed():
	get_node("./QuitMenu/QuitWindow").show()
	get_tree().paused = true

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
	$AudioStreamPlayer.play()
	
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

func _on_multiplayer_timer_timeout():
	var player: CharacterBody2D = get_node_or_null("Player")
	if player != null:
		var sprite = player.get_node("AnimatedSprite2D") as AnimatedSprite2D
		$MultiplayerController.send_player_info(Preferences.get_player_name(), sprite.animation, sprite.frame, player.position, sprite.flip_h)

func _on_quit_window_close_requested():
	get_tree().paused = false
	$QuitMenu/QuitWindow.hide()

func _on_volumeslider_value_changed(value: float):
	$AudioStreamPlayer.set_volume_db(value - 100)

func _on_color_check_toggled(toggled: bool):
	var level = get_node("GeneratedLevel")
	ProjectSettings.set_setting("rendering/environment/defaults/color_blind_mode", toggled)
	if toggled:
		for child in level.get_children():
			if child.get_node("Collider").is_in_group("Red"):
				child.get_node("Sprite").material = ShaderMaterial.new()
				child.get_node("Sprite").material.shader = load("res://colorblind_filter.gdshader")
				child.get_node("Sprite").material.set("shader_parameter/tile_size", 2)
				child.get_node("Sprite").material.set("shader_parameter/pattern", load("res://red_cb.png"))
			if child.get_node("Collider").is_in_group("Green"):
				child.get_node("Sprite").material = ShaderMaterial.new()
				child.get_node("Sprite").material.shader = load("res://colorblind_filter.gdshader")
				child.get_node("Sprite").material.set("shader_parameter/tile_size", 2)
				child.get_node("Sprite").material.set("shader_parameter/pattern", load("res://green_cb.png"))
			if child.get_node("Collider").is_in_group("Blue"):
				child.get_node("Sprite").material = ShaderMaterial.new()
				child.get_node("Sprite").material.shader = load("res://colorblind_filter.gdshader")
				child.get_node("Sprite").material.set("shader_parameter/tile_size", 2)
				child.get_node("Sprite").material.set("shader_parameter/pattern", load("res://blue_cb.png"))
	else:
		for child in level.get_children():
			if child.get_node("Collider").is_in_group("Red"):
				child.get_node("Sprite").set_material(null)
			if child.get_node("Collider").is_in_group("Green"):
				child.get_node("Sprite").set_material(null)
			if child.get_node("Collider").is_in_group("Blue"):
				child.get_node("Sprite").set_material(null)


func _on_dark_check_toggled(toggled: bool):
	if toggled:
		RenderingServer.set_default_clear_color(Color(0, 0, 0, 1))
		$Player.get_node("AnimatedSprite2D").material = ShaderMaterial.new()
		$Player.get_node("AnimatedSprite2D").material.shader = load("res://color_invert.gdshader")
		$QuitButton.material = ShaderMaterial.new()
		$QuitButton.material.shader = load("res://color_invert.gdshader")
		var level = get_node("GeneratedLevel")
		for child in level.get_children():
			if child.get_node("Collider").is_in_group("Black"):
				child.get_node("Sprite").material = ShaderMaterial.new()
				child.get_node("Sprite").material.shader = load("res://color_invert.gdshader")
	else:
		RenderingServer.set_default_clear_color(Color(1, 1, 1, 1))
		$Player.get_node("AnimatedSprite2D").set_material(null)
		$QuitButton.set_material(null)
		var level = get_node("GeneratedLevel")
		for child in level.get_children():
			if child.get_node("Collider").is_in_group("Black"):
				child.get_node("Sprite").set_material(null)


func _on_resume_button_pressed():
	get_tree().paused = false
	$QuitMenu/QuitWindow.hide()
