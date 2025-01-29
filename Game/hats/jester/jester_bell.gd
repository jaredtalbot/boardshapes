extends Node2D

@onready var bell_sprite = $BellSprite

const HALF_PI = PI / 2.0

var last_x = 0
var sway = 0

func _ready():
	last_x = global_position.x

func _physics_process(delta):
	var xv = global_position.x - last_x
	sway = move_toward(sway, xv, delta * 8)
	global_rotation = sin(clampf(sway, -HALF_PI, HALF_PI)) * (HALF_PI)
	last_x = global_position.x
