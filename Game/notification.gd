class_name Notification extends Control

signal finished

func play_animation() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_CENTER_RIGHT)
	var on_screen_position = position
	set_anchors_and_offsets_preset(Control.PRESET_CENTER_LEFT)
	set_anchors_preset(Control.PRESET_CENTER_RIGHT, true)
	var off_screen_position = position
	
	var tween = create_tween()
	tween.tween_property(self, "position", on_screen_position, 0.5) \
		.set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
	tween.tween_property(self, "position", off_screen_position, 0.5) \
		.set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN).set_delay(5.0)
	tween.tween_callback(finished.emit)
	
	_on_play_animation()

func _on_play_animation() -> void:
	pass
