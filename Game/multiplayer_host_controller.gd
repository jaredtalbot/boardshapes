class_name MultiplayerHostController extends Node

var web_server_url: String = ProjectSettings.get_setting("application/boardwalk/web_server_url")

var socket = WebSocketPeer.new()
var current_ready_state: WebSocketPeer.State:
	set(value):
		if current_ready_state != value:
			match value:
				WebSocketPeer.STATE_OPEN:
					host_connected.emit()
				WebSocketPeer.STATE_CLOSED:
					host_disconnected.emit()
		current_ready_state = value
	get():
		return current_ready_state

var lobby_id: String
var players := {}

signal host_connected
signal host_disconnected

func _ready():
	socket.connect_to_url(MPGlobals.to_websocket_url(web_server_url))

func _process(delta):
	socket.poll()
	for k in players:
		players[k].rtc.poll()
	
	current_ready_state = socket.get_ready_state()
	if current_ready_state == WebSocketPeer.STATE_OPEN:
		while socket.get_available_packet_count() > 0:
			var json = JSON.stringify(socket.get_packet().get_string_from_utf8())
			if json is Dictionary:
				match json.get("type"):
					"your_id":
						lobby_id = json.get("id")
					"sdp":
						var sdp = json.content
						if sdp is not Dictionary \
							or sdp.get("type") is not String \
							or sdp.get("sdp") is not String:
							continue
						var lobby_player: LobbyPlayer = players.get(json.playerId)
						if lobby_player == null:
							lobby_player = create_new_lobby_player(json.playerId)
							players[json.playerId] = lobby_player
						lobby_player.rtc.set_remote_description(sdp.type, sdp.sdp)
						lobby_player.rtc.session_description_created \
							.connect(_on_lobby_player_session_description_created.bind(lobby_player))
					"ice":
						var ice = json.content
						if ice is not Dictionary \
							or ice.get("media") is not String \
							or (ice.get("index") is not int and ice.get("index") is not float) \
							or ice.get("name") is not String:
							continue
						var lobby_player: LobbyPlayer = players.get(json.playerId)
						if lobby_player == null:
							lobby_player = create_new_lobby_player(json.playerId)
							players[json.playerId] = lobby_player
						lobby_player.rtc.add_ice_candidate(ice.media, ice.index, ice.name)

func _on_lobby_player_session_description_created(type: String, sdp: String, lobby_player: LobbyPlayer):
	lobby_player.rtc.set_local_description(type, sdp)
	socket.send_text(JSON.stringify({
		"type": "sdp",
		"content": {
			"type": type,
			"sdp": sdp,
		},
		"playerId": lobby_player.player_id,
	}))

func _on_lobby_player_ice_candidate_created(media: String, index: int, name: String, lobby_player: LobbyPlayer):
	lobby_player.rtc.add_ice_candidate(media, index, name)
	socket.send_text(JSON.stringify({
		"type": "ice",
		"content": {
			"media": media,
			"index": index,
			"name": name,
		},
		"playerId": lobby_player.player_id,
	}))

func create_new_lobby_player(player_id: String) -> LobbyPlayer:
	var lobby_player = LobbyPlayer.new()
	lobby_player.rtc = WebRTCPeerConnection.new()
	lobby_player.channels[MPGlobals.MESSAGES] = \
		lobby_player.rtc.create_data_channel("messages", {
			"negotiate": true,
			"id": 0,
		})
	lobby_player.channels[MPGlobals.STATUS] = \
		lobby_player.rtc.create_data_channel("status", {
			"negotiate": true,
			"id": 1,
			"maxPacketLifeTime": 35,
		})
	lobby_player.channels[MPGlobals.RESOURCE] = \
		lobby_player.rtc.create_data_channel("resource", {
			"negotiate": true,
			"id": 2,
		})
	lobby_player.player_id = player_id
	return lobby_player

class LobbyPlayer extends RefCounted:
	var rtc: WebRTCPeerConnection
	var channels: Array[WebRTCDataChannel]
	var player_id: String
