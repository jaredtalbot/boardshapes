extends Sprite2D

@onready var sparkle_particles: CPUParticles2D = $SparkleParticles

func _process(delta):
	var crown_scale = global_scale.length()
	sparkle_particles.scale_amount_max = crown_scale
	sparkle_particles.scale_amount_min = crown_scale
