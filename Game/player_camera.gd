extends Camera2D

#maybe this could be an autoload instead but idrc rn
static var target_zoom := 1.0

func _ready() -> void:
	call_deferred("reset_zoom")

func reset_zoom():
	var final_zoom = get_final_zoom()
	zoom = Vector2(final_zoom, final_zoom)

func _process(delta: float) -> void:
	if Input.is_action_just_pressed("zoom"):
		target_zoom = 1.5 if target_zoom == 1.0 else 1.0
	
	var final_zoom = get_final_zoom()
	zoom = zoom.move_toward(Vector2(final_zoom, final_zoom), \
		maxf(abs(final_zoom - zoom.x) / 2.0, 0.01) * delta * 20)

func get_final_zoom():
	return target_zoom if not get_tree().paused and is_visible_in_tree() else 1.0
