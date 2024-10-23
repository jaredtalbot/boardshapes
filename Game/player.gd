extends CharacterBody2D

@export var SPEED = 500.0
@export var JUMP_VELOCITY = -400.0
@export var acceleration = 1500
@export var wall_slide_speed = 50

@export var dash_speed = 700
var is_dashing = false

var can_jump = false

@export var wall_jump_power = 500

@onready var test_animation = $AnimatedSprite2D
var initial_position : Vector2
	
func _ready():
	test_animation.play("idle animation")
	initial_position = position

func _on_coyote_timer_timeout():
	can_jump = false

var air_time := 0.0

func _physics_process(delta):
	if (position.y > 1092.31863281):
		velocity.x = 0
		velocity.y = 0 
		test_animation.play("death")
		set_physics_process(false)
		var death_timer = get_tree().create_timer(1.0416)
		await death_timer.timeout
		set_physics_process(true)
		position = initial_position
	# Add the gravity.
	if not is_on_floor():
		velocity += get_gravity() * delta
		air_time += delta
		if air_time > 0.25 and test_animation.animation != &"jumping":
			test_animation.play(&"jumping")
			test_animation.set_frame_and_progress(7, 0.0)
		if test_animation.animation == &"jumping":
			if test_animation.frame >= 7 and velocity.y < 0:
				test_animation.set_frame_and_progress(7, 0.0)
	else:
		air_time = 0.0
		can_jump = true
	
	if is_on_floor() == false and can_jump == true and $coyote_timer.is_stopped():
		$coyote_timer.start()
	
	if is_on_floor() and velocity.x == 0:
		test_animation.play("idle animation")
	
	# Handle jump.
	if Input.is_action_just_pressed("jump") and can_jump == true:
		velocity.y = JUMP_VELOCITY
		can_jump = false
		test_animation.play(&"jumping")
		if velocity.x > 0:
			test_animation.flip_h = false
		elif velocity.x < 0:
			test_animation.flip_h = true
		
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
			velocity.x = move_toward(velocity.x, direction * SPEED, acceleration * delta)
			if velocity.x > 0:
				test_animation.flip_h = false
			elif velocity.x < 0:
				test_animation.flip_h = true
			if is_on_floor():
				test_animation.play("running")
	else:
		velocity.x = move_toward(velocity.x, 0, acceleration * delta)

	move_and_slide()
	
	if is_on_wall() and !is_on_floor():
		velocity.y = wall_slide_speed
		test_animation.play("sliding")
		if velocity.x > 0:
			test_animation.flip_h = false
		elif velocity.x < 0:
			test_animation.flip_h = true
	
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


#func _on_visible_on_screen_notifier_2d_screen_exited() -> void:
########		velocity.y = 0 
		#test_animation.play("death")
		#set_physics_process(false)
		#var death_timer = get_tree().create_timer(1.0416)
		#await death_timer.timeout
		#set_physics_process(true)
		#position = initial_position
	
	
	# Replace with function body.
