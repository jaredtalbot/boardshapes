class_name PirateUnlockChecker extends UnlockChecker


static func get_hat_id() -> String:
	return "pirate"

func _pos_set():
	unlock_me()

func _connect_level_signals(level: Level):
	level.loaded.connect(func(): 
		level.pos_set.connect(_pos_set))
