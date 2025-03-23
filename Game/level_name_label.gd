extends SubViewportContainer

var LEVEL_TITLE_FONT = preload("res://The Brownies.otf")

@onready var label: RichTextLabel = %Label

var tween: Tween

## Returns true if the level name was successfully shown, false otherwise.
func display_level_name(level_path: String, queue_free_after = true) -> bool:
	show()
	var level_index = CampaignLevels.all_levels.find_custom(path_equals.bind(level_path))
	if level_index == -1:
		if queue_free_after:
			queue_free()
		else:
			hide()
		return false
	var level: Dictionary = CampaignLevels.all_levels[level_index]
	
	var level_name: String = level.get("name")
	
	if level_name == null or level_name == "":
		if queue_free_after:
			queue_free()
		else:
			hide()
		return false
	
	var level_name_lines = level_name.split("\n", false)
	for i in range(level_name_lines.size()):
		level_name_lines[i] = " " + level_name_lines[i] + "\n"
	
	label.clear()
	label.push_font(LEVEL_TITLE_FONT)
	label.push_color(Color.WHITE)
	label.push_outline_color(Color.BLACK)
	label.push_outline_size(12)
	label.push_font_size(64)
	# title
	label.add_text(level_name_lines[0])
	label.pop()
	label.push_font_size(40)
	for line in level_name_lines.slice(1):
		label.add_text(line)
	label.pop_all()
	
	var mat = (material as ShaderMaterial)
	mat.set_shader_parameter("edge", 0.0)
	mat.set_shader_parameter("reverse", false)
	
	tween = create_tween()
	tween.tween_method(func(x): mat.set_shader_parameter("edge", x), 0.0, 1.0, 3.0).set_delay(0.5)
	tween.tween_callback(mat.set_shader_parameter.bind("reverse", true)).set_delay(4.0)
	tween.tween_method(func(x): mat.set_shader_parameter("edge", x), 1.0, 0.0, 3.0)
	
	if queue_free_after:
		tween.tween_callback(queue_free)
	else:
		tween.tween_callback(hide)
	
	return true

static func path_equals(level_obj: Dictionary, level_path: String):
	return level_obj.get("path") == level_path
