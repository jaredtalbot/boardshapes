extends Control

var hat_name: String
var hat_scene: PackedScene
var hat_description: String
var hat_unlock_hint: String

@onready var hat_holder = $HatHolder

func load_hat_from_json(json: Dictionary):
	hat_scene = load(json.path) if json.get("path") is String else null
	hat_name = json.get("name", "Hat")
	hat_description = json.get("description", "Some sort of hat.")
	hat_unlock_hint = json.get("unlock_hint", "???")
	
	call_deferred("set_hat", hat_scene)

func set_hat(hat: PackedScene):
	assert(hat_holder.get_child_count() < 2)
	if hat != null:
		var new_hat := hat.instantiate()
		new_hat.position = Vector2.ZERO
		if hat_holder.get_child_count() > 0:
			var existing_hat := hat_holder.get_child(0)
			existing_hat.replace_by(new_hat)
			existing_hat.queue_free()
		else:
			hat_holder.add_child(new_hat)
		$NoHatText.hide()
	else:
		if hat_holder.get_child_count() > 0:
			hat_holder.get_child(0).queue_free()
		$NoHatText.show()
