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
var last_position_was_floor: bool
var last_position_was_wall: bool
	
func _ready():
	test_animation.play("idle animation")
	initial_position = position

var air_time := 0.0

func _on_dash_timer_timeout():
	is_dashing = false
	$dash_timer.stop()

func _on_dash_cooldown_timer_timeout():
	$dash_cooldown_timer.stop()

func _on_coyote_timer_timeout():
	$coyote_timer.stop()
	can_jump = false

func _death():
	velocity.x = 0
	velocity.y = 0
	test_animation.play("death")
	set_physics_process(false)
	var death_timer = get_tree().create_timer(1.0416)
	await death_timer.timeout
	set_physics_process(true)
	position = initial_position

func _physics_process(delta):
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
		
	if is_on_floor():
		can_jump = true
		last_position_was_floor = true
		last_position_was_wall = false
	
	if is_on_floor() == false and can_jump == true and $coyote_timer.is_stopped():
		$coyote_timer.start()
	
	
	if is_on_floor() and velocity.x == 0:
		test_animation.play("idle animation")
			
	
	# Handle jump.
	if Input.is_action_just_pressed("jump") and can_jump == true:
		if last_position_was_floor:
			velocity.y = JUMP_VELOCITY
			can_jump = false
			test_animation.play(&"jumping")
			if velocity.x > 0:
				test_animation.flip_h = false
			elif velocity.x < 0:
				test_animation.flip_h = true
		if last_position_was_wall:
			velocity.y = JUMP_VELOCITY
			velocity.x = get_wall_normal().x * wall_jump_power
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
			velocity.x = move_toward(velocity.x, direction * SPEED, acceleration * delta)
			if velocity.x > 0:
				test_animation.flip_h = false
			elif velocity.x < 0:
				test_animation.flip_h = true
			if is_on_floor():
				test_animation.play("running")
	else:
		velocity.x = move_toward(velocity.x, 0, acceleration * delta)

	$AfterImageEmitter.enabled = is_dashing

	move_and_slide()
	
	if is_on_wall() and !is_on_floor():
		last_position_was_floor = false
		last_position_was_wall = true
		can_jump = true
		velocity.y = wall_slide_speed
		test_animation.play("sliding")
		if velocity.x > 0:
			test_animation.flip_h = false
		elif velocity.x < 0:
			test_animation.flip_h = true
	
	if is_on_wall() == false and can_jump == true and $coyote_timer.is_stopped():
		$coyote_timer.start()

	if (position.y > get_viewport_rect().end.y):
		_death()
	
	for index in range(get_slide_collision_count()):
		# We get one of the collisions with the player
		var collision = get_slide_collision(index)

		# If the collision is with ground
		if collision.get_collider() == null:
			continue
		
		if collision.get_collider().is_in_group("Red"):
			_death()
		elif collision.get_collider().is_in_group("Green"):
			if velocity.x > 0:
				velocity.x = SPEED * 2
			elif velocity.x < 0:
				velocity.x = -SPEED * 2
		elif collision.get_collider().is_in_group("Blue"):
			velocity.y = -750
