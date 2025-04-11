class_name Level extends Node

var base_url = ProjectSettings.get_setting("application/boardwalk/web_server_url")

@onready var loading_indicator = $LoadingIndicator
@onready var multiplayer_timer = $MultiplayerTimer
@onready var multiplayer_controller = $MultiplayerController
@onready var level_name_label = %LevelNameLabel

signal loaded
signal started
signal completed
signal pos_set

var player: Player

var level_meta := {}
var current_challenge_level: String = ""
var current_campaign_level: String = ""
## workaround to https://github.com/godotengine/godot/issues/104004
var should_load_next_level = false

func create_level(img: Image, options: Dictionary):
	loading_indicator.show()
	loading_indicator.set_text("Uploading Level...")
	if Preferences.dark_mode:
		loading_indicator.set_text_color(Color(1, 1, 1, 1))
	var buffer = img.save_png_to_buffer()
	var request = FileUploader.upload_buffer(base_url + "/api/create-shapes", buffer, "image.png", HTTPClient.METHOD_POST, "image", options)
	request.request_completed.connect(_on_response_received)

func load_level(level_data: Variant):
	if level_data is PackedByteArray:
		level_data = level_data.get_string_from_utf8()
	
	if level_data is not String:
		show_error("Could not load level.")
		return
	
	var json = JSON.parse_string(level_data)
	
	if json is not Dictionary or json.get("startPos") is not Dictionary \
		or json.get("endPos") is not Dictionary or json.get("regions") is not Array:
		show_error("Could not load level.")
		return
	
	## Convert shapes from old level files
	if (json["regions"] as Array).any(func(x): return (x as Dictionary).has("mesh")):
		for region: Dictionary in json["regions"]:
			if not region.has("mesh"):
				show_error("Could not load level.")
			var mesh = region["mesh"]
			var shape = []
			shape.resize(len(mesh) * 2)
			for i in range(len(mesh)):
				shape[i*2] = mesh[i]["x"]
				shape[i*2+1] = mesh[i]["y"]
			region["shape"] = shape
	
	var generated_level = LevelGenerator.generate_nodes(json["regions"])
	if generated_level == null:
		show_error("Could not load level.")
		return
	
	call_deferred("add_child", generated_level)
	var start_pos = json["startPos"]
	var end_pos = json["endPos"]
	
	level_meta["hash"] = str(hash(level_data))
	level_meta["start_pos"] = Vector2(start_pos["x"], start_pos["y"])
	level_meta["end_pos"] = Vector2(end_pos["x"], end_pos["y"])

func _on_response_received(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray):
	if response_code != HTTPClient.RESPONSE_OK:
		show_error(body, response_code)
		return
	
	var level_data = body.get_string_from_utf8()
	var generated_level = LevelGenerator.generate_nodes(level_data)
	if generated_level == null:
		show_error("Could not generate level with server response.")
		return
	add_child(generated_level)
	level_meta["hash"] = str(hash(level_data))
	initialize_game()

func initialize_game():
	player = add_player()
	var multiplayer_id = level_meta.get("hash")
	if multiplayer_id:
		multiplayer_controller.try_connect(multiplayer_id)
	loading_indicator.hide()
	AccessibilityShaderManager.apply_shaders()
	loaded.emit()
	if not Music.playing:
		Music.play_all_layers()
	
	if current_campaign_level:
		level_name_label.display_level_name(current_campaign_level)
	elif current_challenge_level:
		level_name_label.display_level_name(current_challenge_level)
	
	var start_pos = level_meta.get("start_pos")
	var end_pos = level_meta.get("end_pos")
	
	if not start_pos and not end_pos:
		get_tree().paused = true
		player.hide()
		$StartEndSelection/StartSelect.disabled = false
		$StartEndSelection/StartSelect.show()
		Music.drum_layer.volume_db = linear_to_db(Music.volume/3.0)
		Music.drum_layer.pitch_scale = 0.7
		Music.sample_layer.volume_db = linear_to_db(0.0)
		Music.sample_layer.pitch_scale = 0.7
	elif start_pos and end_pos:
		player.initial_position = start_pos
		player.position = player.initial_position
		var goal = $Goal
		goal.position = end_pos
		goal.show()
		start_game()
	else:
		assert(false, "make sure to set either both start and end positions or neither of them")

func start_game():
	get_tree().paused = false
	$TouchScreenControls.show()
	player.show()
	started.emit()
	var tween = create_tween().set_parallel()
	tween.tween_method(Music.set_layer_volume.bind(Music.drum_layer),
		db_to_linear(Music.drum_layer.volume_db), Music.volume, 1.0)
	tween.tween_method(Music.set_layer_volume.bind(Music.sample_layer),
		db_to_linear(Music.sample_layer.volume_db), Music.volume, 1.0)
	tween.tween_property(Music.drum_layer, "pitch_scale", 1.0, 1.0)
	tween.tween_property(Music.sample_layer, "pitch_scale", 1.0, 1.0)

func show_error(body: Variant, error_code: int = 0):
	var error_dialog = ErrorDialog.new()
	if error_code != 0:
		error_dialog.set_text_to_error_message(body, error_code)
	else:
		error_dialog.set_text_to_error_message(body)
	error_dialog.show()
	add_child(error_dialog)
	error_dialog.confirmed.connect(go_back)
	error_dialog.canceled.connect(go_back)
	

func add_player() -> Player:
	var player = preload("res://player.tscn").instantiate()
	add_child(player)
	return player
	

func _set_player_start():
	player.initial_position = player.get_global_mouse_position()
	player.position = player.initial_position
	player.show()
	$StartEndSelection/StartSelect.disabled = true
	$StartEndSelection/StartSelect.hide()
	$StartEndSelection/EndSelect.disabled = false
	$StartEndSelection/EndSelect.show()

func _set_goal_position():
	var goal = $Goal
	goal.position = player.get_global_mouse_position()
	$StartEndSelection/EndSelect.disabled = true
	$StartEndSelection/EndSelect.hide()
	$Goal.show()
	pos_set.emit()
	start_game()
	
func _goal_reached(_node):
	completed.emit()
	if current_campaign_level != "":
		var currlevel = CampaignLevels.levels.data.map(func(l): return l.path).find(current_campaign_level)
		if currlevel + 1 == len(CampaignLevels.levels.data):
			player.set_physics_process(false)
			$VictoryScreen.show()
			%Restart.call_deferred("grab_focus")
		else:
			player.set_physics_process(false)
			should_load_next_level = true
	else:
		player.set_physics_process(false)
		$VictoryScreen.show()
		%Restart.call_deferred("grab_focus")

func _process(_delta):
	if should_load_next_level:
		go_to_next_campaign_level()
		should_load_next_level = false

## workaround to https://github.com/godotengine/godot/issues/104004
func go_to_next_campaign_level():
	var currlevel = CampaignLevels.levels.data.map(func(l): return l.path).find(current_campaign_level)
	var nextlevel = CampaignLevels.levels.data[currlevel + 1].path
	var transition = ScreenTransitioner.custom_transition()
	var next_level_node = preload("res://level.tscn").instantiate()
	next_level_node.current_campaign_level = nextlevel
	var task_id = WorkerThreadPool.add_task(func():
		var level_data = FileAccess.get_file_as_string(nextlevel)
		next_level_node.load_level(level_data)
	)
	transition.transition_midway.connect(func():
		WorkerThreadPool.wait_for_task_completion(task_id)
		add_sibling(next_level_node)
		get_tree().set_deferred("current_scene", next_level_node)
		next_level_node.call_deferred("initialize_game")
		queue_free()
	)
	transition.transition_canceled.connect(next_level_node.queue_free)

func _on_restart_button_pressed():
	$VictoryScreen.hide()
	$Goal.hide()
	
	loading_indicator.hide()
	get_tree().paused = true
	$StartEndSelection/StartSelect.disabled = false
	$StartEndSelection/StartSelect.show()
	player.set_physics_process(true)

func _on_multiplayer_timer_timeout():
	var player: CharacterBody2D = get_node_or_null("Player")
	var hatPos: Marker2D = get_node_or_null("Player/HatPivot/HatPos") 
	if player != null:
		var sprite = player.get_node("AnimatedSprite2D") as AnimatedSprite2D
		var hatId: String
		if hatPos != null:
			if hatPos.get_child_count() > 0:
				hatId = hatPos.get_child(0).get_meta("hatId")
			else:
				hatId = "nohat"
		$MultiplayerController.send_player_info(Preferences.player_name, sprite.animation, sprite.frame, player.position, hatId, hatPos.position, hatPos.rotation, sprite.flip_h)

func _input(event):
	if event.is_action_pressed("pause"):
		pause()
		get_viewport().set_input_as_handled()

func _on_pause_menu_close_requested():
	unpause()

func _on_resume_button_pressed():
	unpause()

func unpause():
	get_tree().paused = false
	%PauseMenu.hide()
	
func pause():
	get_tree().paused = true
	$PauseMenu.show()

func _on_pause_button_pressed():
	pause()

func _on_back_button_pressed():
	get_tree().paused = false
	go_back()

func go_back():
	ScreenTransitioner.change_scene_to_file("res://menus/start_menu.tscn")

func _on_exit_to_main_menu_button_pressed():
	ScreenTransitioner.change_scene_to_file("res://menus/main.tscn")

func save_level():
	var goal = $Goal as Area2D
	var generated_level = $GeneratedLevel as LevelGenerator.GeneratedLevel
	var level_info = {
		"startPos": { "x": player.initial_position.x, "y": player.initial_position.y },
		"endPos": { "x": goal.position.x, "y": goal.position.y },
		"regions": generated_level.regions
	}
	var json = JSON.stringify(level_info)
	if OS.has_feature("web"):
		JavaScriptBridge.download_buffer(json.to_utf8_buffer(), "level.boardwalk", "application/json")
	else:
		var file := FileAccess.open("user://level.boardwalk", FileAccess.WRITE)
		file.store_string(json)
		file.close()
		OS.shell_show_in_file_manager(ProjectSettings.globalize_path("user://level.boardwalk"))
