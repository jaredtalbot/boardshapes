class_name WebPickFileDialog extends Node
# credit to https://alexduggan1.github.io/Guides/Godot4LetUserUploadFileOnWeb/ for part of the code

@export var accept: PackedStringArray = [".jpeg",".jpg",".png","image/jpeg","image/png"]

static var next_id = 0

static func get_next_id():
	var id = next_id
	next_id += 1
	return id

var _onload = JavaScriptBridge.create_callback(_file_loaded_callback)

signal file_loaded(content: PackedByteArray, filename: String)

var id: int

func _ready():
	if not OS.has_feature("web"):
		queue_free()
		return
	id = get_next_id()
	JavaScriptBridge.eval(("""
var uploadedFile{id} = null;
var uploadedFileType{id} = "";
var callbacks{id} = { onfileload: null }

if(document.getElementById('selectFiles{id}') == null) {
	g = document.createElement('input');
	g.setAttribute("id", "selectFiles{id}");
	g.setAttribute("type", "file");
	g.setAttribute("accept", "%s");
	document.body.append(g);
}

document.getElementById('selectFiles{id}').onchange = async function() {
	let files = document.getElementById('selectFiles{id}').files;
	if (files.length <= 0) {
		return false;
	}
	
	uploadedFile{id} = await files.item(0).arrayBuffer();
	uploadedFileType{id} = files.item(0).name;
	
	callbacks{id}.onfileload();
};
	""" % ",".join(accept)).replace("{id}", str(id)), true)
	var callbacksObj = JavaScriptBridge.get_interface("callbacks" + str(id))
	callbacksObj.onfileload = _onload

func show():
	JavaScriptBridge.eval("document.getElementById('selectFiles%s').click();" % str(id), true)
	
func _file_loaded_callback(_args):
	var content = JavaScriptBridge.eval("uploadedFile" + str(id), true)
	var filename = JavaScriptBridge.eval("uploadedFileType" + str(id), true)
	if content is PackedByteArray:
		file_loaded.emit(content, filename)
	JavaScriptBridge.eval("document.getElementById('selectFiles%s').value = \"\";" % str(id), true)

func _exit_tree():
	JavaScriptBridge.eval("""
var uploadedFile{id} = undefined;
var uploadedFileType{id} = undefined;
var callbacks{id} = undefined;
document.getElementById('selectFiles{id}').remove();
""".replace("{id}", str(id)), true)
	request_ready() # why not
