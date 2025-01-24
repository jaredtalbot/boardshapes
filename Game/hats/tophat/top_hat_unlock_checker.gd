class_name TopHatUnlockChecker extends UnlockChecker

static func get_hat_id() -> String:
	return "tophat"

func _connect_level_signals(level: Level):
	level.completed.connect(Unlocks.unlock_hat.bind("tophat"))
