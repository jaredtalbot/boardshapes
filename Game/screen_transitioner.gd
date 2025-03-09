extends Node

@onready var white_screen = %WhiteScreen

var current_transition: ScreenTransition
var transition_tween: Tween
var next_scene_path: String

func custom_transition() -> ScreenTransition:
	if current_transition:
		cancel_transition()
	var transition = ScreenTransition.new()
	
	var mat = (white_screen.material as ShaderMaterial)
	mat.set_shader_parameter("reverse", false)
	mat.set_shader_parameter("edge", 0.0)
	
	transition_tween = create_tween()
	current_transition = transition
	transition_tween.set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN_OUT)
	transition_tween.tween_callback(white_screen.show)
	transition_tween.tween_callback(transition.transition_started.emit)
	transition_tween.tween_method(func(x): mat.set_shader_parameter("edge", x), 0.0, 1.0, 0.8)
	transition_tween.tween_callback(transition.transition_midway.emit)
	transition_tween.tween_callback(transition.set.bind("past_midway", true))
	transition_tween.tween_callback(mat.set_shader_parameter.bind("reverse", true))
	transition_tween.tween_method(func(x): mat.set_shader_parameter("edge", x), 1.0, 0.0, 0.8)
	transition_tween.tween_callback(transition.transition_ended.emit)
	transition_tween.tween_callback(white_screen.hide)
	transition_tween.tween_callback(set.bind("screen_transitioning", false))
	transition_tween.tween_callback(transition.free)
	
	return transition
	

func change_scene_to_file(path: String):
	next_scene_path = path
	if OS.has_feature("threads"):
		ResourceLoader.load_threaded_request(path)
	if not current_transition or current_transition.past_midway:
		var transition = custom_transition()
		transition.transition_midway.connect(_perform_scene_change)

func _perform_scene_change():
	assert(ResourceLoader.exists(next_scene_path))
	var next_scene: PackedScene
	if OS.has_feature("threads"):
		match ResourceLoader.load_threaded_get_status(next_scene_path):
			ResourceLoader.THREAD_LOAD_IN_PROGRESS, ResourceLoader.THREAD_LOAD_LOADED:
				next_scene = ResourceLoader.load_threaded_get(next_scene_path)
			ResourceLoader.THREAD_LOAD_INVALID_RESOURCE, ResourceLoader.THREAD_LOAD_FAILED:
				cancel_transition()
	else:
		next_scene = load(next_scene_path)
	
	if not next_scene:
		cancel_transition()
	
	get_tree().change_scene_to_packed(next_scene)

func cancel_transition():
	var mat = (white_screen.material as ShaderMaterial)
	if transition_tween:
		transition_tween.kill()
	mat.set_shader_parameter("reverse", false)
	mat.set_shader_parameter("edge", 0.0)
	white_screen.hide()
	if current_transition:
		if current_transition.past_midway:
			current_transition.transition_ended.emit()
		else:
			current_transition.transition_canceled.emit()
		current_transition.free.call_deferred()
		current_transition = null

class ScreenTransition extends Object:
	var past_midway := false
	## Emit when the transition is starting (before the screen is wiped)
	signal transition_started
	## Emits when the screen is fully wiped.
	## Any function that loads the next screen should be connected here.
	signal transition_midway
	## Emits when the transition has fully completed or was "canceled" past the midway point.
	signal transition_ended
	## Emits when the transition was canceled before the midway point.
	signal transition_canceled
	
