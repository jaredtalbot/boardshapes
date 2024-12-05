class_name MultiplayerController extends Node

var multiplayer_server_url: String = ProjectSettings.get_setting("application/boardwalk/multiplayer_server_url")

@export var status_indicator: MPIndicator

var socket = WebSocketPeer.new()

var last_lobby_id: String

func _ready():
	status_indicator.refresh_button.pressed.connect(retry)

func retry():
	try_connect(last_lobby_id)

func try_connect(lobby_id: String):
	last_lobby_id = lobby_id
	var join_url = multiplayer_server_url + "/join"
	if join_url.begins_with("http://"):
		join_url = join_url.replace("http://", "ws://")
	elif join_url.begins_with("https://"):
		join_url = join_url.replace("https://", "wss://")
	socket.connect_to_url(join_url + "?lobby=%s" % lobby_id)

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
	socket.poll()
	var state = socket.get_ready_state()
	
	match state:
		WebSocketPeer.STATE_CONNECTING:
			status_indicator.set_status(MPIndicator.MPConnectionStatus.CONNECTING)
		WebSocketPeer.STATE_OPEN:
			status_indicator.set_status(MPIndicator.MPConnectionStatus.CONNECTED)
			update_players()
		WebSocketPeer.STATE_CLOSED:
			status_indicator.set_status(MPIndicator.MPConnectionStatus.DISCONNECTED)

func update_players():
	for i in range(socket.get_available_packet_count()):
		var json_string = socket.get_packet().get_string_from_utf8()
		var json_obj = JSON.parse_string(json_string)
		if json_obj is Dictionary \
			and json_obj.get("id") is String \
			and json_obj.get("name") is String \
			and json_obj.get("animation") is String \
			and json_obj.get("frame") is float or json_obj.get("frame") is int \
			and json_obj.get("position") is Dictionary \
			and json_obj["position"].get("x") is float or json_obj["position"].get("x") is int \
			and json_obj["position"].get("y") is float or json_obj["position"].get("y") is int \
			and json_obj.get("facingLeft") is bool:
			var ghost: AnimatedSprite2D
			ghost = get_node_or_null(json_obj["id"])
			if ghost == null:
				ghost = preload("res://ghost_player.tscn").instantiate()
				add_child(ghost)
				ghost.name = json_obj["id"]
			ghost.set_player_tag(json_obj["name"])
			ghost.animation = json_obj["animation"]
			ghost.frame = json_obj["frame"]
			ghost.position = Vector2(json_obj["position"]["x"], json_obj["position"]["y"])
			ghost.flip_h = json_obj["facingLeft"]
			ghost.last_updated = Time.get_unix_time_from_system()

func send_player_info(name: String, animation: String, frame: int, position: Vector2, facingLeft: bool):
	if socket.get_ready_state() == WebSocketPeer.STATE_OPEN:
		var info_dict = {
			"name": name,
			"animation": animation,
			"frame": frame,
			"position": {"x": position.x, "y": position.y},
			"facingLeft": facingLeft
		}
		var json = JSON.stringify(info_dict)
		socket.send_text(json)

func _exit_tree():
	socket.close()
