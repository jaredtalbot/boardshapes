extends Notification

@onready var hat_name_label = %HatNameLabel
@onready var hat_unlock_hint_label = %HatUnlockHintLabel
@onready var hat_holder = %HatHolder
@onready var unlock_sound = $UnlockSound
@onready var confetti = %Confetti

var hat_scene: PackedScene

func load_hat_by_id(id: String) -> void:
	var json: Dictionary
	for v in Unlocks.HAT_LIST.data:
		if v.id == id:
			json = v
	hat_scene = load(json.path) if json.get("path") is String else null
	hat_name_label.text = json.get("name", "Hat")
	hat_unlock_hint_label.text = json.get("unlock_hint", "???")
	
	call_deferred("set_hat", hat_scene)

func _on_play_animation() -> void:
	unlock_sound.play()
	confetti.emitting = true

func set_hat(hat: PackedScene) -> void:
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
