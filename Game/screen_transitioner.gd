extends Node

@onready var white_screen = %WhiteScreen

var screen_transitioning := false
var next_scene_path: String

func custom_transition() -> CustomScreenTransition:
	if screen_transitioning:
		return null
	var transition = CustomScreenTransition.new()
	
	var mat = (white_screen.material as ShaderMaterial)
	mat.set_shader_parameter("reverse", false)
	mat.set_shader_parameter("edge", 0.0)
	
	var tween = create_tween()
	screen_transitioning = true
	tween.set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN_OUT)
	tween.tween_callback(transition.transition_started.emit)
	tween.tween_method(func(x): mat.set_shader_parameter("edge", x), 0.0, 1.0, 1.0)
	tween.tween_callback(transition.mid_transition.emit)
	tween.tween_callback(mat.set_shader_parameter.bind("reverse", true))
	tween.tween_method(func(x): mat.set_shader_parameter("edge", x), 1.0, 0.0, 1.0)
	tween.tween_callback(transition.transition_ended.emit)
	tween.tween_callback(white_screen.hide)
	tween.tween_callback(set.bind("screen_transitioning", false))
	tween.tween_callback(transition.free)
	
	return transition
	

func change_scene_to_file(path: String):
	next_scene_path = path
	if not screen_transitioning:
		screen_transitioning = true
		if OS.has_feature("threads"):
			ResourceLoader.load_threaded_request(path)
		white_screen.show()
		
		var mat = (white_screen.material as ShaderMaterial)
		mat.set_shader_parameter("reverse", false)
		mat.set_shader_parameter("edge", 0.0)
		
		var tween = create_tween()
		tween.set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN_OUT)
		tween.tween_method(func(x): mat.set_shader_parameter("edge", x), 0.0, 1.0, 1.0)
		tween.tween_callback(_perform_scene_change)

func _perform_scene_change():
	assert(ResourceLoader.exists(next_scene_path))
	var next_scene: PackedScene
	if OS.has_feature("threads"):
		match ResourceLoader.load_threaded_get_status(next_scene_path):
			ResourceLoader.THREAD_LOAD_IN_PROGRESS, ResourceLoader.THREAD_LOAD_LOADED:
				next_scene = ResourceLoader.load_threaded_get(next_scene_path)
			ResourceLoader.THREAD_LOAD_INVALID_RESOURCE, ResourceLoader.THREAD_LOAD_FAILED:
				# abandon scene change
				screen_transitioning = false
				white_screen.hide()
				return
	else:
		next_scene = load(next_scene_path)
	
	if not next_scene:
		# abandon scene change
		screen_transitioning = false
		white_screen.hide()
		return
	
	get_tree().change_scene_to_packed(next_scene)
	
	var mat = (white_screen.material as ShaderMaterial)
	mat.set_shader_parameter("reverse", true)
	mat.set_shader_parameter("edge", 1.0)
	
	var tween = create_tween()
	tween.set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN_OUT)
	tween.tween_method(func(x): mat.set_shader_parameter("edge", x), 1.0, 0.0, 1.0)
	tween.tween_callback(white_screen.hide)
	tween.tween_callback(set.bind("screen_transitioning", false))

class CustomScreenTransition extends Object:
	signal transition_started
	signal mid_transition
	signal transition_ended
