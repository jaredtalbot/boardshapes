class_name MultiplayerController extends Node

const HAT_LIST = preload("res://hats/hat_list.json")

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
func _process(_delta):
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
			and json_obj.get("hatId") is String \
			and json_obj["hatPosition"].get("x") is float or json_obj["hatPosition"].get("x") is int \
			and json_obj["hatPosition"].get("y") is float or json_obj["hatPosition"].get("y") is int \
			and json_obj.get("hatRotation") is float or json_obj.get("hatRotation") is int \
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
			if json_obj["hatId"] != "nohat" && ghost.get_node("HatPivot/HatPos").get_child_count() < 1:
				for hat_json in HAT_LIST.data:
					if hat_json["id"] == json_obj["hatId"]:
						var hat_scene = load(hat_json["path"]) if hat_json.get("path") is String else null
						ghost.get_node("HatPivot/HatPos").add_child(hat_scene.instantiate())
			ghost.get_node("HatPivot/HatPos").position = Vector2(json_obj["hatPosition"]["x"], json_obj["hatPosition"]["y"])
			ghost.get_node("HatPivot/HatPos").rotation = json_obj["hatRotation"] 
			ghost.flip_h = json_obj["facingLeft"]
			if json_obj["facingLeft"]:
				ghost.get_node("HatPivot").scale.x = -7.813
			else:
				ghost.get_node("HatPivot").scale.x = 7.813
			ghost.last_updated = Time.get_unix_time_from_system()

func send_player_info(name: String, animation: String, frame: int, position: Vector2, hatId: String, hatPosition: Vector2, hatRotation: float, facingLeft: bool):
	if socket.get_ready_state() == WebSocketPeer.STATE_OPEN:
		var info_dict = {
			"name": name,
			"animation": animation,
			"frame": frame,
			"position": {"x": position.x, "y": position.y},
			"hatId": hatId,
			"hatPosition": {"x": hatPosition.x, "y": hatPosition.y},
			"hatRotation": hatRotation,
			"facingLeft": facingLeft
		}
		var json = JSON.stringify(info_dict)
		socket.send_text(json)

func _exit_tree():
	socket.close()
