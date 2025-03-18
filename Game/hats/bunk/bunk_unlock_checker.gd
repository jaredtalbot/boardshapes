class_name BunkUnlockChecker extends UnlockChecker

const BUNK_META = "link:https://shenangelos.itch.io/bunk"

static func get_hat_id() -> String:
	return "bunk"

func _ready():
	GlobalSignals.customize_menu_loaded.connect(_on_customize_menu_loaded)

func _on_customize_menu_loaded(customize_menu: CustomizeMenu):
	var unlock_hint_label: RichTextLabel = customize_menu.get_node("HatInfoDisplay/HatUnlockHintLabel")
	
	unlock_hint_label.meta_clicked.connect(_on_unlock_hint_label_meta_clicked)

func _on_unlock_hint_label_meta_clicked(meta: Variant):
	if meta == BUNK_META:
		unlock_me()
