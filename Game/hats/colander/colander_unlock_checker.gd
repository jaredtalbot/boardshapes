class_name ColanderUnlockChecker extends UnlockChecker

static func get_hat_id() -> String:
	return "colander"

func _connect_level_signals(level: Level):
	level.completed.connect(_check_if_last_level.bind(level))

func _check_if_last_level(level: Level):
	if level.current_challenge_level == CampaignLevels.challenge_levels.data[-1].path:
		unlock_me()
