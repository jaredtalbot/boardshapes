class_name MPGlobals extends Object

enum {MESSAGES = 0, STATUS = 1, RESOURCE = 2}

static func to_websocket_url(url: String):
	var ws_url = url
	if ws_url.begins_with("http://"):
		ws_url = ws_url.replace("http://", "ws://")
	elif ws_url.begins_with("https://"):
		ws_url = ws_url.replace("https://", "wss://")
	return ws_url
