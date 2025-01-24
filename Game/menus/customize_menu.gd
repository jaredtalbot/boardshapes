extends Node

const HAT_LIST = preload("res://hats/hat_list.json")
const HAT_PREVIEW = preload("res://hats/hat_preview.tscn")

@onready var cosmetic_preview = $PreviewAnchor/CosmeticPreview
@onready var hat_select = %HatSelect
@onready var multiplayer_name_field = %MultiplayerNameField


func _ready():
	Music.stop_all_layers()
	multiplayer_name_field.text = Preferences.player_name
	$HatInfoDisplay.hide()
	cosmetic_preview.equip_hat(Preferences.hat_scene)
	Preferences.hat_scene_changed.connect(cosmetic_preview.equip_hat)
	
	for hat_json in HAT_LIST.data:
		var hat_preview = HAT_PREVIEW.instantiate()
		hat_preview.load_hat_from_json(hat_json)
		hat_select.add_child(hat_preview)
		hat_preview.pressed.connect(set_hat.bind(hat_preview.hat_scene, hat_preview.unlocked))
		hat_preview.focus_entered.connect(set_hat_info_display.bind(hat_preview.hat_name, hat_preview.hat_description, hat_preview.hat_unlock_hint))

func set_hat(new_hat: PackedScene, unlocked: bool):
	if unlocked:
		Preferences.hat_scene = new_hat

func set_hat_info_display(hat_name: String, hat_description: String, hat_unlock_hint: String):
	$HatInfoDisplay/HatNameLabel.text = hat_name
	$HatInfoDisplay/HatDescriptionLabel.text = hat_description
	$HatInfoDisplay/HatUnlockHintLabel.text = hat_unlock_hint
	if not hat_unlock_hint or hat_unlock_hint == "???":
		$HatInfoDisplay/UnlockHintSeparator.hide()
		$HatInfoDisplay/HatUnlockHintLabel.hide()
	else:
		$HatInfoDisplay/UnlockHintSeparator.show()
		$HatInfoDisplay/HatUnlockHintLabel.show()
	$HatInfoDisplay.show()

func _on_back_button_pressed() -> void:
	get_tree().change_scene_to_file("res://menus/start_menu.tscn")

func _on_multiplayer_name_field_text_changed(new_text):
	Preferences.player_name = new_text
