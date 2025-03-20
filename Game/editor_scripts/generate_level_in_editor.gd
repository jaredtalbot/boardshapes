@tool
extends EditorScript

func _run():
	var fd = preload("res://pick_boardwalk_file_dialog.tscn").instantiate()
	
	fd.file_selected.connect(_on_file_selected.bind(fd))
	fd.canceled.connect(fd.queue_free)
	EditorInterface.popup_dialog(fd)
	## need ref to not get freed :/
	fd.set_meta("ref", self)

func _on_file_selected(path: String, fd: FileDialog):
	fd.queue_free()
	
	var data = FileAccess.get_file_as_string(path)
	
	var parser = JSON.new()
	
	var err = parser.parse(data)
	
	if err:
		printerr("invalid boardwalk file, could not parse json: %s" % parser.get_error_message())
		return
	
	if parser.data is not Dictionary:
		printerr("invalid boardwalk file, not an object")
		return
	
	var json = parser.data as Dictionary
	var regions = json.get("regions")
	
	if regions is not Array:
		printerr("invalid boardwalk file, did not find regions")
	
	var root = get_scene()
	
	if not root:
		add_root_node(Node.new())
		root = get_scene()
	
	var level = LevelGenerator.generate_nodes(regions)
	
	print(level)
	
	root.add_child(level)
