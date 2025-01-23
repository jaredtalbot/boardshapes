class_name UnlockChecker extends Node

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	UnlockCheckerManager.new_level_added.connect(_connect_level_signals)

func _connect_level_signals(level: Level):
	pass

static func get_hat_id() -> String:
	assert(false, "get_hat_id should be overridden")
	return ""
