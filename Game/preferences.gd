extends Node

var player_name := "":
	get():
		return player_name if player_name != "" else "Player"

var hat_scene: PackedScene:
	set(value):
		hat_scene = value
		hat_scene_changed.emit(hat_scene)

signal hat_scene_changed(hat_scene: PackedScene)
