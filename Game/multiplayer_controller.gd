class_name MultiplayerController extends Node

var web_server_url: String = ProjectSettings.get_setting("application/boardwalk/web_server_url")

var socket = WebSocketPeer.new()

func try_connect(lobby_id: String):
	var join_url = web_server_url + "/api/join"
	if join_url.begins_with("http://"):
		join_url = join_url.replace("http://", "ws://")
	elif join_url.begins_with("https://"):
		join_url = join_url.replace("https://", "wss://")
	socket.connect_to_url(join_url + "?lobby=%s" % lobby_id)

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
	socket.poll()
	var state = socket.get_ready_state()
	
	if state == WebSocketPeer.STATE_OPEN:
		update_players()

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
