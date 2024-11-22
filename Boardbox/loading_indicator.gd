extends VBoxContainer

@onready var loading_icon = $LoadingIconHolder/LoadingIcon
@onready var loading_indicator_text = $LoadingIndicatorText

func _process(delta):
	if visible:
		loading_icon.rotation = loading_icon.rotation + fmod(delta * TAU, 360)
	else:
		loading_icon.rotation = 0

func set_text(new_text: String):
	loading_indicator_text.text = new_text
