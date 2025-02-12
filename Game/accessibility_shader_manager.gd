extends Node

func _ready():
	Preferences.dark_mode_changed.connect(apply_dark_mode)
	Preferences.colorblind_mode_changed.connect(apply_colorblind_mode)
	apply_shaders()

func apply_shaders() -> void:
	apply_dark_mode(Preferences.dark_mode)
	apply_colorblind_mode(Preferences.colorblind_mode)

func apply_shaders_to_node(node: CanvasItem) -> void:
	apply_dark_mode_to_node(Preferences.dark_mode, node)
	apply_colorblind_mode_to_node(Preferences.colorblind_mode, node)

func apply_dark_mode(enabled: bool) -> void:
	if enabled:
		RenderingServer.set_default_clear_color(Color.BLACK)
		for node in get_tree().get_nodes_in_group("DarkModeInvertColors"):
			if not try_apply_special_dark_mode(node):
				node.material = preload("res://shaders/materials/dark_mode.tres")
	else:
		RenderingServer.set_default_clear_color(Color.WHITE)
		for node in get_tree().get_nodes_in_group("DarkModeInvertColors") \
			+ get_tree().get_nodes_in_group("DarkModeTurnWhite"):
			node.material = null

func apply_dark_mode_to_node(enabled: bool, node: CanvasItem) -> void:
	if node.is_in_group("DarkModeInvertColors"):
		if enabled:
			if not try_apply_special_dark_mode(node):
				node.material = preload("res://shaders/materials/dark_mode.tres")
		else:
			node.material = null

## This exists thanks to Godot applying shaders to text controls weirdly.
##
## Returns true if the node was one of those special cases, false otherwise.
func try_apply_special_dark_mode(node: CanvasItem) -> bool:
	if node is Label:
		node.material = ShaderMaterial.new()
		node.material.shader = preload("res://shaders/silhouette.gdshader")
		if node.label_settings:
			node.material.set_shader_parameter("color", node.label_settings.font_color.inverted())
		else:
			node.material \
				.set_shader_parameter("color", node.get_theme_color("font_color", "Label").inverted())
		return true
	elif node is Control and node.has_theme_color("font_color"):
		node.material = ShaderMaterial.new()
		node.material.shader = preload("res://shaders/silhouette.gdshader")
		node.material.set_shader_parameter("color", node.get_theme_color("font_color").inverted())
		return true
	else:
		return false

func apply_colorblind_mode(enabled: bool) -> void:
	var tree := get_tree()
	if enabled:
		for node in tree.get_nodes_in_group("Red"):
			node.material = preload("res://shaders/materials/colorblind_red.tres")
		for node in tree.get_nodes_in_group("Green"):
			node.material = preload("res://shaders/materials/colorblind_green.tres")
		for node in tree.get_nodes_in_group("Blue"):
			node.material = preload("res://shaders/materials/colorblind_blue.tres")
	else:
		for node in tree.get_nodes_in_group("Red") \
			+ tree.get_nodes_in_group("Green") \
			+ tree.get_nodes_in_group("Blue"):
			node.material = null

func apply_colorblind_mode_to_node(enabled: bool, node: CanvasItem) -> void:
	if enabled:
		if node.is_in_group("Red"):
			node.material = preload("res://shaders/materials/colorblind_red.tres")
		if node.is_in_group("Green"):
			node.material = preload("res://shaders/materials/colorblind_green.tres")
		if node.is_in_group("Blue"):
			node.material = preload("res://shaders/materials/colorblind_blue.tres")
	else:
		if node.is_in_group("Red") or node.is_in_group("Green") or node.is_in_group("Blue"):
			node.material = null
