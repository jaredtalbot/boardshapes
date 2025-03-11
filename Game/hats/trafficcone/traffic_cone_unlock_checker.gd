class_name TrafficConeUnlockChecker extends UnlockChecker

var death_count: int = 0

static func get_hat_id() -> String:
	return "trafficcone"

func _died():
	death_count += 1
	_check_unlock()

func _check_unlock():
	if death_count >= 100:
		Unlocks.unlock_hat(get_hat_id())

func _connect_level_signals(level: Level):
	level.loaded.connect(func():
		level.player.died.connect(_died))
