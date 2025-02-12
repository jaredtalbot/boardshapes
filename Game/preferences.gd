extends Node

var player_name := "":
	get():
		return player_name if player_name != "" else "Player"

var hat_scene: PackedScene = null:
	set(value):
		hat_scene = value
		hat_scene_changed.emit(hat_scene)

var volume := 0.5:
	set(value):
		volume = clampf(value, 0.0, 1.0)
		volume_changed.emit(volume)

var touchscreen_button_scale := 5.0:
	set(value):
		touchscreen_button_scale = clampf(value, 1.0, 6.0)
		touchscreen_button_scale_changed.emit(touchscreen_button_scale)

var dark_mode := false:
	set(value):
		dark_mode = value
		dark_mode_changed.emit(dark_mode)

var colorblind_mode := false:
	set(value):
		colorblind_mode = value
		colorblind_mode_changed.emit(colorblind_mode)

var save_timer: Timer

signal hat_scene_changed(hat_scene: PackedScene)
signal touchscreen_button_scale_changed(new_scale: float)
signal volume_changed(new_volume: float)
signal dark_mode_changed(new_value: bool)
signal colorblind_mode_changed(new_value: bool)
signal saved

func _ready():
	save_timer = Timer.new()
	save_timer.one_shot = true
	save_timer.wait_time = 1.0
	save_timer.process_mode = Node.PROCESS_MODE_ALWAYS
	save_timer.timeout.connect(_on_save_timer_timeout)
	add_child(save_timer)
	var err = load_preferences()
	if err:
		Notifications.show_message("Failed to load preferences:\n" + error_string(err))

func load_preferences() -> Error:
	var file := FileAccess.open("user://preferences.json", FileAccess.READ)
	var err := FileAccess.get_open_error()
	if err != OK:
		if err == ERR_FILE_NOT_FOUND:
			save_preferences()
			return OK
		return err
	
	var json = JSON.parse_string(file.get_as_text())
	if json is Dictionary:
		player_name = json.get("player_name", player_name)
		var hat_id = json.get("hat_id")
		if hat_id:
			for hat in Unlocks.HAT_LIST.data:
				if hat.get("id") == hat_id:
					hat_scene = load(hat.path) if hat.get("path") else null
					break
		volume = json.get("volume", volume)
		touchscreen_button_scale = json.get("touchscreen_button_scale", touchscreen_button_scale)
		dark_mode = json.get("dark_mode", dark_mode)
		colorblind_mode = json.get("colorblind_mode", colorblind_mode)
	else:
		save_preferences()
	return OK

func save_when_ready() -> void:
	save_timer.start()

func _on_save_timer_timeout() -> void:
	save_preferences()

func save_preferences() -> Error:
	var file := FileAccess.open("user://preferences.json", FileAccess.WRITE)
	var err := FileAccess.get_open_error()
	if err:
		return err
	
	var hat_id := ""
	if hat_scene:
		for hat in Unlocks.HAT_LIST.data:
			if hat_scene.resource_path == hat.get("path"):
				hat_id = hat.get("id", "")
				break
	file.store_string(JSON.stringify({
		"player_name": player_name,
		"hat_id": hat_id,
		"volume": volume,
		"touchscreen_button_scale": touchscreen_button_scale,
		"dark_mode": dark_mode,
		"colorblind_mode": colorblind_mode,
	}))
	saved.emit()
	return OK
