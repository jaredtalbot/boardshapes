extends RichTextLabel

var LEVEL_TITLE_FONT = preload("res://The Brownies.otf")

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
	
	clear()
	push_font(LEVEL_TITLE_FONT)
	push_color(Color.WHITE)
	push_outline_color(Color.BLACK)
	push_outline_size(12)
	push_font_size(64)
	# title
	add_text(level_name_lines[0])
	pop()
	push_font_size(20)
	for line in level_name_lines.slice(1):
		add_text(line)
	pop_all()
	
	modulate = Color.TRANSPARENT
	
	tween = create_tween()
	tween.tween_property(self, "modulate", Color.WHITE, 1.0)
	tween.tween_property(self, "modulate", Color.TRANSPARENT, 1.0).set_delay(3)
	
	if queue_free_after:
		tween.tween_callback(queue_free)
	
	return true

static func path_equals(level_obj: Dictionary, level_path: String):
	return level_obj.get("path") == level_path
