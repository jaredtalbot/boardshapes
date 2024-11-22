extends Control

@export var level_scene: PackedScene

@onready var pick_image_file_dialog = $PickImageFileDialog
@onready var web_pick_image_file = $WebPickImageFile
@onready var image_confirmation = $ImageConfirmation

#todo: add non-web file picker

func _on_upload_image_button_pressed():
	if OS.has_feature("web"):
		web_pick_image_file.show()
	else:
		pick_image_file_dialog.show()

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
