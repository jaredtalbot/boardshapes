extends Node

var player_name := "":
	get():
		return player_name if player_name != "" else "Player"

var hat_scene: PackedScene:
	set(value):
		hat_scene = value
		hat_scene_changed.emit(hat_scene)

signal hat_scene_changed(hat_scene: PackedScene)
signal saved

func _ready():
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
		player_name = json.get("player_name", "")
		var hat_id = json.get("hat_id", "")
		for hat in Unlocks.HAT_LIST.data:
			if hat.get("id") == hat_id:
				hat_scene = load(hat.path) if hat.get("path") else null
				break
	else:
		save_preferences()
	return OK

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
	}))
	saved.emit()
	return OK
