extends Node

@onready var hat_holder = %Hat
@onready var container = $SubViewportContainer

var hat_scene: PackedScene
var hat_list: Array = []

func _ready():
	load_hats()

func load_hats():
	for hat_data in Unlocks.HAT_LIST.data:
		var json = hat_data
		hat_scene = load(json.path) if json.get("path") is String else null
		
		if hat_scene:
			var new_subviewport = create_subviewport()

			set_hat_for_subviewport(new_subviewport, hat_scene, json.id)
			
			var new_particles = create_particles(new_subviewport)

			hat_list.append({
				"subviewport": new_subviewport,
				"particles": new_particles,
				"hat": hat_scene
			})

func create_subviewport() -> SubViewport:
	var viewport = SubViewport.new()
	viewport.set_size(Vector2(512, 512))
	viewport.transparent_bg = true

	add_child(viewport)

	return viewport

func set_hat_for_subviewport(viewport: SubViewport, hat: PackedScene, id: String):
	var new_hat = hat.instantiate()
	new_hat.scale *= 10
	new_hat.position = Vector2(256,256)
	_update_modulation(new_hat, id)
	Unlocks.updated.connect(func(): _update_modulation(new_hat, id))
	viewport.add_child(new_hat)

func _update_modulation(new_hat: Node2D, id: String):
	if id in Unlocks.unlocked_hat_ids:
		new_hat.modulate = Color8(255, 255, 255, 125)
	else:
		new_hat.modulate = Color8(0, 0, 0, 125)
	
func create_particles(viewport: SubViewport):
	var particles_scene = preload("res://sample_emitter.tscn")
	var new_particles = particles_scene.instantiate()
	new_particles.texture = viewport.get_texture()
	new_particles.preprocess = randf_range(0, 10)
	#new_particles.process_material.gravity.x = randf_range(-25, -10)
	add_child(new_particles)
	return new_particles
