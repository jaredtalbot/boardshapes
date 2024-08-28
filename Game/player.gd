extends CharacterBody2D

@export var SPEED = 300.0
@export var JUMP_VELOCITY = -400.0
@export var acceleration = 25

@export var dash_speed = 700
var is_dashing = false

var can_jump = false

@export var wall_jump_power = 500

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
		
	if $dash_cooldown_timer.time_left == 0:
		_on_dash_cooldown_timer_timeout()
		
	if Input.is_action_just_pressed("Dash") and is_dashing == false and $dash_timer.is_stopped() and $dash_cooldown_timer.time_left == 0:
			$dash_timer.start()
			$dash_cooldown_timer.start()
			is_dashing = true
	if direction:
		if is_dashing == true:
			velocity.y = 0;
			velocity.x = direction * dash_speed
			
		else:
			velocity.x = move_toward(velocity.x, direction * SPEED, acceleration)
	else:
		velocity.x = move_toward(velocity.x, 0, acceleration)

	move_and_slide()
	
	if is_on_wall() and !is_on_floor():
		velocity.y = 50
	
	if is_on_wall() and Input.is_action_pressed("jump"):
		$wall_timer.start()
		velocity.y = JUMP_VELOCITY
		velocity.x = get_wall_normal().x * wall_jump_power
		can_jump = false
		

func _on_dash_timer_timeout():
	is_dashing = false
	$dash_timer.stop()


func _on_dash_cooldown_timer_timeout():
	$dash_cooldown_timer.stop()
