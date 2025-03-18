extends Node

signal new_level_added(level: Level)
signal customize_menu_loaded(customize_menu: CustomizeMenu)

func _ready():
	get_tree().node_added.connect(_on_tree_node_added)

func _on_tree_node_added(node: Node):
	if node is Level:
		new_level_added.emit(node)
	elif node is CustomizeMenu:
		customize_menu_loaded.emit(node)
