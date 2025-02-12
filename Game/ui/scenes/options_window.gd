extends Window

func _on_close_button_pressed():
	hide()

func _on_delete_save_data_button_pressed():
	# todo: confirmation
	var err = Unlocks.clear_unlocks()
	if err == OK:
		Notifications.show_message("Save data deleted.")
	else:
		Notifications.show_message("Error with deleting save data:\n" + error_string(err))
