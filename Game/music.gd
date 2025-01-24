extends Node

var playing:
	get:
		return drum_layer.playing or sample_layer.playing
var volume = 0.5

@onready var drum_layer: AudioStreamPlayer = $DrumLayer
@onready var sample_layer: AudioStreamPlayer = $SampleLayer

static func set_layer_volume(new_volume: float, layer: AudioStreamPlayer):
	layer.volume_db = linear_to_db(new_volume)

func set_volume(new_volume: float, update_layers: bool = true) -> void:
	volume = new_volume
	if update_layers:
		restore_layers_volume()

func restore_layers_volume() -> void:
	set_layer_volume(volume, drum_layer)
	set_layer_volume(volume, sample_layer)

func stop_all_layers() -> void:
	drum_layer.stop()
	sample_layer.stop()

func play_all_layers(from_position: float = 0.0, restore_defaults: bool = true) -> void:
	if restore_defaults:
		restore_layers_volume()
		drum_layer.pitch_scale = 1.0
		sample_layer.pitch_scale = 1.0
	drum_layer.play(from_position)
	sample_layer.play(from_position)
