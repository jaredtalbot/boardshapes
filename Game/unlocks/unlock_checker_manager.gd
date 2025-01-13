extends Node

signal new_level_added(level: Level)

func _ready() -> void:
	get_tree().node_added.connect(_on_tree_node_added)
	var unlock_checkers = ProjectSettings.get_global_class_list() \
		.filter(func(x): return x["base"] == &"UnlockChecker")
	for uc in unlock_checkers:
		var script: GDScript = load(uc["path"])
		var hat_id = script.get_hat_id()
		if hat_id not in Unlocks.unlocked_hat_ids:
			add_child(script.new())

func _on_tree_node_added(node: Node):
	if node is Level:
		new_level_added.emit(node)
