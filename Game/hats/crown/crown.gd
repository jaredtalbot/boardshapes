extends Sprite2D

@onready var sparkle_particles = $SparkleParticles

func _ready():
	var process_material: ParticleProcessMaterial = sparkle_particles.process_material
	var crown_scale = global_scale.length()
	process_material.scale_max = crown_scale
	process_material.scale_min = crown_scale
