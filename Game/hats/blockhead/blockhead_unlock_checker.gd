class_name BlockheadUnlockChecker extends UnlockChecker

const UNLOCK_AREA_SCENE = preload("res://hats/blockhead/blockhead_unlock_area.tscn")

static func get_hat_id() -> String:
	return "blockhead"

func _connect_level_signals(level: Level):
	if level.current_campaign_level == "res://campaign/j'slvl.boardwalk":
		level.loaded.connect(_add_unlock_area.bind(level))

func _add_unlock_area(level: Level):
	var unlock_area_node = UNLOCK_AREA_SCENE.instantiate()
	level.add_child(unlock_area_node)
	
	var area = unlock_area_node.get_node("Area2D") as Area2D
	area.body_entered.connect(_area_entered)

func _area_entered(node: Node):
	if node is Player:
		unlock_me()
