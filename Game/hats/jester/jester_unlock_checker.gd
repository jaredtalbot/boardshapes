class_name JesterUnlockChecker extends UnlockChecker

static func get_hat_id() -> String:
	return "jester"

func _ready():
	if OS.is_debug_build():
		unlock_me()
		return
	var valid = JavaScriptBridge.eval("""
const url = new URL(window.location);

const today = new Date();
today.setHours(0, 0, 0, 0);
const todayHash = today.toISOString().split("").reduce((a, b) => {
  a = ((a << 5) - a) + b.charCodeAt(0);
  return a & a;
}, 0);

url.searchParams.get("unlock") === `${todayHash}`;
""")
	if valid:
		unlock_me()
		JavaScriptBridge.eval("""
const url = new URL(window.location);
url.search = "";
history.replaceState(null, undefined, url);
""")
