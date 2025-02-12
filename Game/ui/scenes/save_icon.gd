extends TextureRect

var tween: Tween

func _ready():
	Unlocks.saved.connect(appear)
	Preferences.saved.connect(appear)

func appear():
	if tween:
		tween.kill()
	tween = create_tween()
	modulate = Color(Color.WHITE, 0.8)
	tween.tween_property(self, "modulate", Color.TRANSPARENT, 1.0)
