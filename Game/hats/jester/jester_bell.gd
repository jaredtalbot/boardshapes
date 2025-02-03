extends Node2D

const HALF_PI = PI / 2.0

var last_x = 0
var av = 0

func _ready():
	last_x = global_position.x

func _physics_process(delta):
	scale.y = sign(scale.y) * sign(global_scale.y)
	var xv = global_position.x - last_x
	var force = Vector2(xv * delta * 5, 1)
	
	av += Vector2.from_angle(global_rotation + PI).dot(force) * delta
	av *= 1 - delta
	
	global_rotation = clampf(global_rotation + av, -HALF_PI, HALF_PI)
	last_x = global_position.x
