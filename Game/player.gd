extends CharacterBody2D

@export var SPEED = 500.0
@export var JUMP_VELOCITY = -400.0
@export var acceleration = 1500
@export var wall_slide_speed = 50

@export var dash_speed = 900

var can_jump = false

@export var wall_jump_power = 500

@onready var test_animation = $AnimatedSprite2D
var initial_position : Vector2
var last_position_was_floor: bool
var last_position_was_wall: bool
var touched_green: bool = false
var bonked_wall: bool = false
var dash_direction: int
	
func _ready():
	test_animation.play("idle animation")
	initial_position = position

var air_time := 0.0

func _on_coyote_timer_timeout():
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
	if velocity.x > 0:
		test_animation.flip_h = false
	elif velocity.x < 0:
		test_animation.flip_h = true
	
	var is_dashing = not $dash_timer.is_stopped() or touched_green
	bonked_wall = false
	
	# Add the gravity.
	if not is_on_floor():
		velocity += get_gravity() * delta
		air_time += delta
		if air_time > 0.05 and test_animation.animation != &"jumping":
			test_animation.play(&"jumping")
		if test_animation.animation == &"jumping":
			if test_animation.frame >= 11 and velocity.y < 0:
				test_animation.set_frame_and_progress(11, 0.0)
	else:
		air_time = 0.0
		can_jump = true
		
	if is_on_floor():
		$land_particles.emitting = true
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
		if last_position_was_wall:
			velocity.y = JUMP_VELOCITY
			velocity.x = get_wall_normal().x * wall_jump_power
			can_jump = false
			
		
	# Get the input direction and handle the movement/deceleration.
	# As good practice, you should replace UI actions with custom gameplay actions.
	var direction = Input.get_axis("left", "right")
		
	if Input.is_action_just_pressed("Dash") and not is_dashing \
		and $dash_cooldown_timer.is_stopped() and not is_on_wall():
			$dash_timer.start()
			$dash_cooldown_timer.start()
			dash_direction = get_direction()
			
	if is_dashing:
			if not test_move(transform, Vector2(0, 2)):
				velocity.y = 0
			velocity.x = dash_direction * dash_speed
			test_animation.play("dash")
	else:
		if direction:
			velocity.x = move_toward(velocity.x, direction * SPEED, acceleration * delta)
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
		$dash_timer.stop()
		test_animation.play("sliding")
		$slide_particles_right.emitting = true
		$slide_particles_left.emitting = true
		if test_animation.flip_h == false:
			$slide_particles_left.emitting = false
		elif test_animation.flip_h == true:
			$slide_particles_right.emitting = false
	
	if !is_on_wall():
		$slide_particles_right.emitting = false
		$slide_particles_left.emitting = false
	
	if is_dashing and is_on_floor() and is_on_wall():
		bonked_wall = true
		velocity = Vector2(get_wall_normal().x, -0.5) * 500
		
	
	if is_on_wall() == false and can_jump == true and $coyote_timer.is_stopped():
		$coyote_timer.start()

	if (position.y > get_viewport_rect().end.y):
		_death()
		
	touched_green = false
	
	for index in range(get_slide_collision_count()):
		# We get one of the collisions with the player
		var collision = get_slide_collision(index)
		
		var collider = collision.get_collider()
		# If the collision is with ground
		if collider == null:
			continue
		
		if collider.is_in_group("Red"):
			_death()
		elif collider.is_in_group("Green") and collision.get_normal().dot(Vector2.UP) > 0.5 and not bonked_wall:
			touched_green = true
			dash_direction = get_direction()
		elif collider.is_in_group("Blue"):
			velocity.y = -750

func get_direction() -> int:
	return -1 if test_animation.flip_h else 1
