extends PanelContainer

@onready var web_server_list = $WebServerList
var button_group = ButtonGroup.new()

var web_server_urls: PackedStringArray = ProjectSettings.get_setting("application/boardwalk/available_web_server_urls")

func _ready():
	if not OS.has_feature("debug"):
		queue_free()
		return
	
	if web_server_urls == null:
		$WebServerList.queue_free()
	
	for url in web_server_urls:
		var url_option = WebServerOption.new()
		url_option.url = url
		url_option.button_group = button_group
		url_option.url_checked.connect(_on_option_url_checked.bind(url_option))
		$WebServerList.add_child(url_option)
	
	(func(): if get_child_count() == 0: queue_free()).call_deferred()

func _on_option_url_checked(valid: bool, option: WebServerOption):
	var current_option = button_group.get_pressed_button()
	if valid and web_server_urls.has(option.url) and (current_option == null \
		or not web_server_urls.has(current_option.url) \
		or web_server_urls.find(option.url) < web_server_urls.find(current_option.url)):
			option.button_pressed = true

class WebServerOption extends CheckBox:
	signal url_checked(valid: bool)
	
	var url: String:
		set(value):
			url = value
			text = value
		get(): return url
	
	var check_request: HTTPRequest
	
	func _ready():
		check_request = HTTPRequest.new()
		add_child(check_request)
		check_request.name = "CheckRequest"
		check_request.request_completed.connect(_on_check_request_completed)
		check_request.request(url, PackedStringArray(), HTTPClient.METHOD_HEAD)
		disabled = true
		icon = preload("res://icons/pending.png")
	
	func _on_check_request_completed(result: int, response_code: int, _headers: PackedStringArray, _body: PackedByteArray):
		var response_code_string = str(response_code)
		if result != HTTPRequest.RESULT_SUCCESS \
			or not (response_code_string.begins_with("2") or response_code_string.begins_with("3")):
			icon = preload("res://icons/x.png")
			url_checked.emit(false)
		else:
			icon = preload("res://icons/check.png")
			disabled = false
			url_checked.emit(true)
		check_request.queue_free()
	
	func _toggled(toggled_on):
		if toggled_on:
			ProjectSettings.set_setting("application/boardwalk/web_server_url", url)
