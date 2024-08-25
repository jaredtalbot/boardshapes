extends CharacterBody2D

const SPEED = 300.0
const JUMP_VELOCITY = -400.0

var dash_speed = 700
var is_dashing = false

var can_jump = false

func _on_coyote_timer_timeout():
	can_jump = false


func _physics_process(delta):
	
	# Add the gravity.
	if not is_on_floor():
		velocity += get_gravity() * delta
		
	if can_jump == false and is_on_floor() == true:
		can_jump = true
	
	if is_on_floor() == false and can_jump == true and $coyote_timer.is_stopped():
		$coyote_timer.start()
	
	# Handle jump.
	if Input.is_action_just_pressed("jump") and can_jump == true:
		velocity.y = JUMP_VELOCITY
		can_jump = false
		
	# Get the input direction and handle the movement/deceleration.
	# As good practice, you should replace UI actions with custom gameplay actions.
	var direction = Input.get_axis("left", "right")
	
	if $dash_timer.time_left == 0:
		_on_dash_timer_timeout()
	if Input.is_action_just_pressed("Dash") and is_dashing == false and $dash_timer.is_stopped():
			$dash_timer.start()
			is_dashing = true
	if direction:
		if is_dashing == true:
			velocity.x = direction * dash_speed
			
		else:
			velocity.x = direction * SPEED
	else:
		velocity.x = move_toward(velocity.x, 0, SPEED)

	move_and_slide()

func _on_dash_timer_timeout():
	is_dashing = false
	$dash_timer.stop()
