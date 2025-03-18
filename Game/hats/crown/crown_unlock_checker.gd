class_name CrownUnlockChecker extends UnlockChecker

static func get_hat_id() -> String:
	return "crown"

func _connect_level_signals(level: Level):
	level.completed.connect(_check_if_last_level.bind(level))

func _check_if_last_level(level: Level):
	if level.current_campaign_level == CampaignLevels.levels.data[-1].path:
		unlock_me()
