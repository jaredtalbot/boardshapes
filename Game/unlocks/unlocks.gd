extends Node

const HAT_LIST = preload("res://hats/hat_list.json")

signal updated

var unlocked_hat_ids := PackedStringArray()

func _ready():
	load_unlocks()

func unlock_hat(hat_id: String):
	if hat_id not in unlocked_hat_ids:
		unlocked_hat_ids.append(hat_id)
		updated.emit()
		save_unlocks()

func load_unlocks():
	var file := FileAccess.open("user://unlocks.json", FileAccess.READ)
	var err := FileAccess.get_open_error()
	if err != OK:
		if err == ERR_FILE_NOT_FOUND:
			add_always_unlocked_hats()
			save_unlocks()
			return
		#add some kind of alert to the user that unlocks can't load
		return 
	
	var json = JSON.parse_string(file.get_as_text())
	if json is Dictionary and json.get("hats") is Array:
		unlocked_hat_ids = PackedStringArray(json.hats)
		add_always_unlocked_hats()
	else:
		save_unlocks()

func add_always_unlocked_hats():
	unlocked_hat_ids.append_array(HAT_LIST.data \
		.filter(func(h): return h.get("always_unlocked", false)) \
		.map(func(h): return h.id))

func save_unlocks():
	var file := FileAccess.open("user://unlocks.json", FileAccess.WRITE)
	var err := FileAccess.get_open_error()
	if err:
		return #add some kind of alert to the user that unlocks can't save
	file.store_string(JSON.stringify({
		"hats": unlocked_hat_ids
	}))
