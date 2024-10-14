extends Control

@export var level_scene: PackedScene

@onready var web_pick_image_file = $WebPickImageFile
@onready var image_confirmation = $ImageConfirmation

#todo: add non-web file picker

func _on_upload_image_button_pressed():
	web_pick_image_file.show()
	
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
	level.create_level(img)
	queue_free()
