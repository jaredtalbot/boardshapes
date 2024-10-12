class_name WebPickImageFile extends Node
# credit to https://alexduggan1.github.io/Guides/Godot4LetUserUploadFileOnWeb/ for part of the code

var _onload = JavaScriptBridge.create_callback(_file_loaded_callback)

signal file_loaded(content: PackedByteArray, filename: String)

func _ready():
	JavaScriptBridge.eval("""
var uploadedFile = null;
var uploadedFileType = "";
var callbacks = { onfileload: null }

if(document.getElementById('selectFiles') == null) {
	g = document.createElement('input');
	g.setAttribute("id", "selectFiles");
	g.setAttribute("type", "file");
	g.setAttribute("accept", ".jpeg,.jpg,.png,image/jpeg,image/png");
	document.body.append(g);
}

document.getElementById('selectFiles').onchange = async function() {
	var files = document.getElementById('selectFiles').files;
	if (files.length <= 0) {
		return false;
	}
	
	uploadedFile = await files.item(0).arrayBuffer();
	uploadedFileType = files.item(0).name;
	
	callbacks.onfileload();
};
	""", true)
	var callbacksObj = JavaScriptBridge.get_interface("callbacks")
	callbacksObj.onfileload = _onload

func show():
	JavaScriptBridge.eval("document.getElementById('selectFiles').click();", true)
	
func _file_loaded_callback(args):
	var content = JavaScriptBridge.eval("uploadedFile", true)
	var filename = JavaScriptBridge.eval("uploadedFileType", true)
	if content is PackedByteArray:
		file_loaded.emit(content, filename)
