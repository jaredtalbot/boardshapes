extends Node

var unlocked_hat_paths := PackedStringArray()

func _ready():
	load_unlocks()

func load_unlocks():
	var file := FileAccess.open("user://unlocks.json", FileAccess.READ)
	var err := FileAccess.get_open_error()
	if err != OK:
		if err != ERR_FILE_NOT_FOUND:
			pass #add some kind of alert to the user that unlocks can't load
		save_unlocks()
		return 
	
	var json = JSON.parse_string(file.get_as_text())
	if json is Dictionary and json.get("hats") is Array:
		unlocked_hat_paths = PackedStringArray(json.hats)
	else:
		save_unlocks()

func save_unlocks():
	var file := FileAccess.open("user://unlocks.json", FileAccess.WRITE)
	var err := FileAccess.get_open_error()
	if err:
		return #add some kind of alert to the user that unlocks can't save
	file.store_string(JSON.stringify({
		"hats": unlocked_hat_paths
	}))
