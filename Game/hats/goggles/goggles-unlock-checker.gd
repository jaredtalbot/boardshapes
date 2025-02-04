class_name GogglesUnlockChecker extends UnlockChecker

var has_jumped_or_dashed = false

static func get_hat_id() -> String:
	return "goggles"

func _jumped_or_dashed():
	has_jumped_or_dashed = true

func _c_level_complete(level: Level):
	if has_jumped_or_dashed == false and level.current_campaign_level == "res://campaign/c'slvl.boardwalk":
		Unlocks.unlock_hat("goggles")

func _connect_level_signals(level: Level):
	level.loaded.connect(func(): level.player.jumped.connect(_jumped_or_dashed); level.player.dashed.connect(_jumped_or_dashed))
	has_jumped_or_dashed = false
	level.completed.connect(_c_level_complete.bind(level))
