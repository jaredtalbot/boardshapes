extends CanvasLayer

const HatUnlockNotification = preload("res://hat_unlock_notification.tscn")
const MessageNotification = preload("res://message_notification.tscn")

enum { MESSAGE, HAT_UNLOCK }

var notification_queue: Array[Dictionary] = []
var current_notification: Node

func _ready():
	layer = 3

func _process(delta):
	if len(notification_queue) and (
			not current_notification
			or not is_instance_valid(current_notification)
			or current_notification.is_queued_for_deletion()
		):
		_show_next_notification()

func _show_next_notification():
	var notif_info: Dictionary = notification_queue.pop_front()
	match notif_info.get("type"):
		MESSAGE:
			current_notification = MessageNotification.instantiate()
			add_child(current_notification)
			current_notification.set_message_text(notif_info["message_text"])
		HAT_UNLOCK:
			current_notification = HatUnlockNotification.instantiate()
			add_child(current_notification)
			current_notification.load_hat_by_id(notif_info["hat_id"])
	if current_notification:
		var notif = current_notification
		notif.finished.connect(func(): set_deferred("current_notification", null); notif.queue_free())
		notif.play_animation()

func show_message(message_text: String):
	notification_queue.append({
		"type": MESSAGE,
		"message_text": message_text
	})

func show_hat_unlock(hat_id: String):
	notification_queue.append({
		"type": HAT_UNLOCK,
		"hat_id": hat_id
	})
