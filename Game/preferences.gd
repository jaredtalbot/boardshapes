extends Node

var player_name := ""

func get_player_name():
	return player_name if player_name != "" else "Player"

func set_player_name(new_player_name: String):
	player_name = new_player_name
