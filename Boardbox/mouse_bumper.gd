extends AnimatableBody2D

@onready var collision_circle = $CollisionCircle

var bumping: bool

func _physics_process(delta):
	position = get_window().get_mouse_position()
	bumping =  Input.is_mouse_button_pressed(MOUSE_BUTTON_LEFT)
	collision_circle.set_deferred("disabled", not bumping)
	queue_redraw()

func _draw():
	if (bumping):
		draw_circle(Vector2.ZERO, 10, Color.RED)
