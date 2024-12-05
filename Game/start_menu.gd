extends Control

@export var level_scene: PackedScene

@onready var upload_image_button = $VBoxContainer/HBoxContainer/UploadImageButton
@onready var load_level_button = $VBoxContainer/HBoxContainer/LoadLevelButton
@onready var back_button = $VBoxContainer/BackButton
@onready var pick_image_file_dialog = $PickImageFileDialog
@onready var web_pick_image_file = $WebPickImageFile
@onready var web_pick_boardwalk_file = $WebPickBoardwalkFile
@onready var pick_boardwalk_file_dialog = $PickBoardwalkFileDialog
@onready var image_confirmation = $ImageConfirmation

func _ready():
	if RenderingServer.get_default_clear_color() == Color(0, 0, 0, 1):
		upload_image_button.material = ShaderMaterial.new()
		upload_image_button.material.shader = load("res://color_invert.gdshader")
		load_level_button.material = ShaderMaterial.new()
		load_level_button.material.shader = load("res://color_invert.gdshader")
		back_button.material = ShaderMaterial.new()
		back_button.material.shader = load("res://color_invert.gdshader")
	upload_image_button.call_deferred("grab_focus")

func _on_upload_image_button_pressed():
	if OS.has_feature("web"):
		web_pick_image_file.show()
	else:
		pick_image_file_dialog.show()

func _on_load_level_button_pressed():
	if OS.has_feature("web"):
		web_pick_boardwalk_file.show()
	else:
		pick_boardwalk_file_dialog.show()

func _on_pick_image_file_dialog_file_selected(path):
	var img = Image.load_from_file(path)
	if img == null:
		return
	image_confirmation.set_image(img)
	image_confirmation.show()

func _on_back_button_pressed():
	get_tree().change_scene_to_file("res://main.tscn")

func _on_web_pick_image_file_file_loaded(content: PackedByteArray, filename: String):
	var image = Image.new()
	var error: Error
	if filename.ends_with("png"):
		error = image.load_png_from_buffer(content)
	elif filename.ends_with("jpg") or filename.ends_with("jpeg"):
		error = image.load_jpg_from_buffer(content)
	
	if error != 0:
		return
	
	image_confirmation.set_image(image)
	image_confirmation.show()

func _on_image_confirmation_confirmed():
	var img = image_confirmation.get_image()
	if img == null:
		return
	var level = level_scene.instantiate()
	add_sibling(level)
	level.create_level(img, image_confirmation.get_settings())
	get_tree().set_deferred("current_scene", level)
	queue_free()

func _on_multiplayer_name_field_text_changed(new_text):
	Preferences.set_player_name(new_text)

func _on_web_pick_boardwalk_file_file_loaded(content: PackedByteArray, filename: String):
	load_level_from_data(content)

func _on_pick_boardwalk_file_dialog_file_selected(path):
	var data := FileAccess.get_file_as_string(path)
	if FileAccess.get_open_error() != 0:
		return
	load_level_from_data(data)

func load_level_from_data(level_data):
	var level = level_scene.instantiate()
	add_sibling(level)
	level.load_level(level_data)
	get_tree().set_deferred("current_scene", level)
	queue_free()
