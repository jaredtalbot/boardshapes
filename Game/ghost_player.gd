extends AnimatedSprite2D

var last_updated = 0

@onready var player_tag_panel: PanelContainer = %PlayerTagPanel
@onready var player_tag = %PlayerTag

func _ready():
	AccessibilityShaderManager.apply_shaders_to_node(self)

func _process(_delta):
	var time_since_last_update = Time.get_unix_time_from_system() - last_updated
	self_modulate = Color(Color.WHITE, clampf(3 - time_since_last_update, 0, 0.5))
	player_tag_panel.modulate = Color(Color.WHITE, clampf(3 - time_since_last_update, 0, 1))
	if time_since_last_update > 5:
		queue_free()

func set_player_tag(name: String):
	player_tag.text = name
