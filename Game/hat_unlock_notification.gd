extends TextureRect

@onready var hat_name_label = %HatNameLabel
@onready var hat_unlock_hint_label = %HatUnlockHintLabel
@onready var hat_holder = %HatHolder
@onready var unlock_sound = $UnlockSound
@onready var confetti = %Confetti

var hat_scene: PackedScene

func load_hat_by_id(id: String):
	var json: Dictionary
	for v in Unlocks.HAT_LIST.data:
		if v.id == id:
			json = v
	hat_scene = load(json.path) if json.get("path") is String else null
	hat_name_label.text = json.get("name", "Hat")
	hat_unlock_hint_label.text = json.get("unlock_hint", "???")
	
	call_deferred("set_hat", hat_scene)

func play_animation():
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
	unlock_sound.play()
	confetti.emitting = true

func set_hat(hat: PackedScene):
	assert(hat_holder.get_child_count() < 2)
	if hat != null:
		var new_hat := hat.instantiate()
		new_hat.position = Vector2.ZERO
		if hat_holder.get_child_count() > 0:
			var existing_hat := hat_holder.get_child(0)
			hat_holder.add_child(new_hat)
			existing_hat.queue_free()
		else:
			hat_holder.add_child(new_hat)
	else:
		if hat_holder.get_child_count() > 0:
			hat_holder.get_child(0).queue_free()
