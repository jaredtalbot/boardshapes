extends ConfirmationDialog

@onready var texture_rect = $VBoxContainer/TextureRect
@onready var preserve_color_check = %PreserveColorCheck
@onready var no_color_separation_check = %NoColorSeparationCheck
@onready var allow_white_check = %AllowWhiteCheck


func _ready():
	texture_rect.texture = null

func set_image(image):
	if image is Image:
		image = ImageTexture.create_from_image(image)
	
	if image is ImageTexture:
		texture_rect.texture = image

func get_image() -> Image:
	return (texture_rect.texture as ImageTexture).get_image()

func get_settings():
	return {
		"preserveColor": str(preserve_color_check.button_pressed),
		"noColorSeparation": str(no_color_separation_check.button_pressed),
		"allowWhite": str(allow_white_check.button_pressed)
	}
