class_name MPIndicator extends PanelContainer

enum MPConnectionStatus { CONNECTING = 0, CONNECTED = 1, DISCONNECTED = 2 }

@onready var status_icon = $StatusContainer/StatusIcon
@onready var status_label = $StatusContainer/StatusLabel
@onready var refresh_button = $StatusContainer/RefreshButton as TextureButton

var prev_status: MPConnectionStatus

var connecting_string = "Connecting"
var connected_string = "Connected"
var disconnected_string = "Disconnected"

var connecting_icon = preload("res://icons/pending.png")
var connected_icon = preload("res://icons/check.png")
var disconnected_icon = preload("res://icons/x.png")

func set_status(new_status: MPConnectionStatus):
	if new_status == prev_status:
		return
	prev_status = new_status
	
	match new_status:
		MPConnectionStatus.CONNECTING:
			status_icon.texture = connecting_icon
			status_label.text = connecting_string
			refresh_button.hide()
		MPConnectionStatus.CONNECTED:
			status_icon.texture = connected_icon
			status_label.text = connected_string
			refresh_button.hide()
		MPConnectionStatus.DISCONNECTED:
			status_icon.texture = disconnected_icon
			status_label.text = disconnected_string
			refresh_button.show()
		
