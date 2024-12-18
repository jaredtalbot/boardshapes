extends Node

const HAT_LIST = preload("res://hats/hat_list.json")
const HAT_PREVIEW = preload("res://hats/hat_preview.tscn")

@onready var cosmetic_preview = $PreviewAnchor/CosmeticPreview
@onready var hat_select = %HatSelect

func _ready():
	cosmetic_preview.equip_hat(Preferences.hat_scene)
	Preferences.hat_scene_changed.connect(cosmetic_preview.equip_hat)
	
	for hat_json in HAT_LIST.data:
		var hat_preview = HAT_PREVIEW.instantiate()
		hat_preview.load_hat_from_json(hat_json)
		hat_select.add_child(hat_preview)
		hat_preview.pressed.connect(set_hat.bind(hat_preview.hat_scene))

func set_hat(new_hat: PackedScene):
	Preferences.hat_scene = new_hat

func _on_back_button_pressed() -> void:
	get_tree().change_scene_to_file("res://menus/start_menu.tscn")
