extends CanvasLayer

const BUTTON_SIZE = 64

func _ready():
	update_layout(Preferences.touchscreen_button_scale)
	Preferences.touchscreen_button_scale_changed.connect(update_layout)

func update_layout(new_scale: float):
	$BottomLeft/BottomLeftControlRow/LeftButtonControl.custom_minimum_size = Vector2(BUTTON_SIZE*new_scale, BUTTON_SIZE*new_scale)
	$BottomLeft/BottomLeftControlRow/RightButtonControl.custom_minimum_size = Vector2(BUTTON_SIZE*new_scale, BUTTON_SIZE*new_scale)
	$BottomLeft/BottomLeftControlRow/LeftButtonControl/LeftButton.scale = Vector2(new_scale, new_scale)
	$BottomLeft/BottomLeftControlRow/RightButtonControl/RightButton.scale = Vector2(new_scale, new_scale)
	$BottomRight/BottomRightControlRow.custom_minimum_size = Vector2(0, BUTTON_SIZE*new_scale*1.6)
	$BottomRight/BottomRightControlRow/JumpButtonControl.custom_minimum_size = Vector2(BUTTON_SIZE*new_scale, BUTTON_SIZE*new_scale)
	$BottomRight/BottomRightControlRow/DashButtonControl.custom_minimum_size = Vector2(BUTTON_SIZE*new_scale, BUTTON_SIZE*new_scale)
	$BottomRight/BottomRightControlRow/JumpButtonControl/JumpButton.scale = Vector2(new_scale, new_scale)
	$BottomRight/BottomRightControlRow/DashButtonControl/DashButton.scale = Vector2(new_scale, new_scale)
