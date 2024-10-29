extends ConfirmationDialog

@onready var texture_rect = $VBoxContainer/TextureRect

func _ready():
	texture_rect.texture = ImageTexture.create_from_image(Image.create_empty(0, 0, false, Image.FORMAT_RGBA8))

func set_image(image):
	if image is Image:
		image = ImageTexture.create_from_image(image)
	
	if image is ImageTexture:
		texture_rect.texture = image

func get_image() -> Image:
	return (texture_rect.texture as ImageTexture).get_image()
