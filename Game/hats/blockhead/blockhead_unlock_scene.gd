extends Node

@onready var collect_area := $CollectArea
@onready var blockhead_collectable := $BlockheadCollectable
var level: Level
var bound_player: Player = null

func _ready():
	collect_area.body_entered.connect(bind_to_player)
	level = get_parent()
	while level is not Level:
		level = level.get_parent()
	
	level.completed.connect(try_unlock)

func _process(delta):
	var t = fmod(Time.get_unix_time_from_system(), 2)
	blockhead_collectable.rotation = t * PI * 2
	if bound_player:
		blockhead_collectable.global_position = blockhead_collectable.global_position.move_toward(
			bound_player.global_position + Vector2(cos(t*PI) * 40, sin(t*PI) * 15), delta * 1000)
		blockhead_collectable.z_index = 1 if t > 0.5 and t < 1.5 else -1
	else:
		blockhead_collectable.global_position = blockhead_collectable.global_position.move_toward(
			collect_area.global_position, delta * 1000)
		blockhead_collectable.z_index = 1

func bind_to_player(node: Node):
	if node is Player:
		bound_player = node
		bound_player.died.connect(unbind_from_player)

func unbind_from_player():
	bound_player = null

func try_unlock():
	if bound_player:
		Unlocks.unlock_hat("blockhead")
