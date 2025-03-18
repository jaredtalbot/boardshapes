extends Node

func _ready() -> void:
	reload_checkers()

func reload_checkers():
	var children = get_children()
	for node in children:
		remove_child(node)
		node.queue_free()
	var unlock_checkers = ProjectSettings.get_global_class_list() \
		.filter(func(x): return x["base"] == &"UnlockChecker")
	for uc in unlock_checkers:
		var script: GDScript = load(uc["path"])
		var hat_id = script.get_hat_id()
		if hat_id not in Unlocks.unlocked_hat_ids:
			add_child(script.new())
