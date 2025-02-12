extends Node

## Automatically change game settings when options are tweaked. If turned off, the controls
## will not affect game settings and will only emit signals.
@export var default_behavior = true
## Removes the touchscreen button scale option outside of mobile platforms.
@export var button_scale_mobile_only = true

@onready var volume_slider = $VolumeField/Slider
@onready var touchscreen_button_scale_slider = $TouchscreenButtonScaleField/Slider
@onready var dark_mode_check = $DarkModeCheck
@onready var colorblind_mode_check = $ColorblindModeCheck

signal volume_changed(new_volume: float)
signal touchscreen_button_scale_changed(new_scale: float)
signal dark_mode_changed(new_value: bool)
signal colorblind_mode_changed(new_value: bool)

func _ready():
	volume_slider.set_value_no_signal(Preferences.volume * 100.0)
	touchscreen_button_scale_slider.set_value_no_signal(Preferences.touchscreen_button_scale)
	dark_mode_check.set_pressed_no_signal(Preferences.dark_mode)
	colorblind_mode_check.set_pressed_no_signal(Preferences.colorblind_mode)
	if button_scale_mobile_only and not OS.has_feature("mobile"):
		$TouchscreenButtonScaleField.hide()

func _on_dark_mode_toggled(toggled):
	if default_behavior:
		Preferences.dark_mode = toggled
		Preferences.save_when_ready()
	dark_mode_changed.emit(toggled)

func _on_colorblind_mode_toggled(toggled):
	if default_behavior:
		Preferences.colorblind_mode = toggled
		Preferences.save_when_ready()
	colorblind_mode_changed.emit(toggled)

func _on_volume_slider_value_changed(value: float):
	if default_behavior:
		Music.set_volume(value / 100.0)
	volume_changed.emit(value)

func _on_touchscreen_button_scale_slider_value_changed(value: float):
	if default_behavior:
		Preferences.touchscreen_button_scale = value
		Preferences.save_when_ready()
	touchscreen_button_scale_changed.emit(value)
