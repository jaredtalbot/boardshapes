class_name BunkUnlockChecker extends UnlockChecker

static func get_hat_id() -> String:
	return "bunk"

func _connect_level_signals(level: Level):
	level.completed.connect(Unlocks.unlock_hat.bind("tophat"))
