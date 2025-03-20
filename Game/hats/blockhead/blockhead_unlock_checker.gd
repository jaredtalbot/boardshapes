class_name BlockheadUnlockChecker extends UnlockChecker

const UNLOCK_AREA_SCENE = preload("res://hats/blockhead/blockhead_unlock_scene.tscn")

static func get_hat_id() -> String:
	return "blockhead"

func _connect_level_signals(level: Level):
	if level.current_campaign_level == "res://campaign/j'slvl.boardwalk":
		level.loaded.connect(_add_unlock_area.bind(level))

func _add_unlock_area(level: Level):
	var unlock_area_node = UNLOCK_AREA_SCENE.instantiate()
	level.add_child(unlock_area_node)
