class_name TrafficConeUnlockChecker extends UnlockChecker

var death_count: int = 0
const SAVE_FILE_PATH: String = "user://achievementprogress.cfg"

static func get_hat_id() -> String:
	return "trafficcone"

func _died():
	death_count += 1
	_save_death_count()
	_check_unlock()

func _check_unlock():
	if death_count >= 100:
		Unlocks.unlock_hat(get_hat_id())

func _connect_level_signals(level: Level):
	level.loaded.connect(func():
		_load_death_count()
		if level.player:
			level.player.died.connect(_died)
	)

func _save_death_count():
	var config = ConfigFile.new()
	config.set_value("death_count", "count", death_count)
	var err = config.save(SAVE_FILE_PATH)

func _load_death_count():
	var config = ConfigFile.new()
	var err = config.load(SAVE_FILE_PATH)
	if err == OK:
		death_count = config.get_value("death_count", "count", 0)
		
func clear_death_count():
	death_count = 0
	_save_death_count()
	print("Death count cleared")  
