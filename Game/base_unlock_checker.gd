class_name UnlockChecker extends Node

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	GlobalSignals.new_level_added.connect(_connect_level_signals)

func _connect_level_signals(_level: Level):
	pass

static func get_hat_id() -> String:
	assert(false, "get_hat_id should be overridden")
	return ""

func unlock_me():
	Unlocks.unlock_hat(get_hat_id())
