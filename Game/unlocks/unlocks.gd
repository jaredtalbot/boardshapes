extends Node

const HAT_LIST = preload("res://hats/hat_list.json")

signal updated

var unlocked_hat_ids := PackedStringArray()

func _ready():
	var err = load_unlocks()
	if err:
		Notifications.show_message("Failed to load unlocks:\n" + error_string(err))

func unlock_hat(hat_id: String):
	if hat_id not in unlocked_hat_ids and HAT_LIST.data.any(func(x): return x.id == hat_id):
		unlocked_hat_ids.append(hat_id)
		Notifications.show_hat_unlock(hat_id)
		updated.emit()
		var err = save_unlocks()
		if err:
			Notifications.show_message("Failed to save unlocks:\n" + error_string(err))

func load_unlocks() -> Error:
	var file := FileAccess.open("user://unlocks.json", FileAccess.READ)
	var err := FileAccess.get_open_error()
	if err != OK:
		if err == ERR_FILE_NOT_FOUND:
			add_always_unlocked_hats()
			save_unlocks()
			return OK
		return err
	
	var json = JSON.parse_string(file.get_as_text())
	if json is Dictionary and json.get("hats") is Array:
		unlocked_hat_ids = PackedStringArray(json.hats)
		add_always_unlocked_hats()
	else:
		save_unlocks()
	return OK

func add_always_unlocked_hats():
	unlocked_hat_ids.append_array(HAT_LIST.data \
		.filter(func(h): return h.get("always_unlocked", false) and h.id not in unlocked_hat_ids) \
		.map(func(h): return h.id))

func save_unlocks() -> Error:
	var file := FileAccess.open("user://unlocks.json", FileAccess.WRITE)
	var err := FileAccess.get_open_error()
	if err:
		return err
	file.store_string(JSON.stringify({
		"hats": unlocked_hat_ids
	}))
	return OK

func clear_unlocks() -> Error:
	unlocked_hat_ids = PackedStringArray()
	add_always_unlocked_hats()
	updated.emit()
	return save_unlocks()
