extends Notification

@onready var message_label = $MarginContainer/MessageLabel

func set_message_text(text: String):
	message_label.text = text
