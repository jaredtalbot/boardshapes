extends Node

var crypto = Crypto.new()

const filler_chars = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
var nfiller_chars = filler_chars.length()

func upload_file(url: String, filepath: String, method: HTTPClient.Method, fieldname: String = "file") -> HTTPRequest:
	var boundary = generate_boundary(16)
	
	var file_data = FileAccess.get_file_as_bytes(filepath)
	
	if file_data.size() <= 0:
		print(FileAccess.get_open_error())
		return null
		
	var filename = filepath.get_file()
	var extension = filepath.get_extension()
	var mime_type: String
	match extension:
		"jpg", "jpeg":
			mime_type = "image/jpeg"
		"png":
			mime_type = "image/png"
		_:
			mime_type = "application/octet-stream"
	
	var body = PackedByteArray()
	
	body.append_array("--".to_utf8_buffer())
	body.append_array(boundary)
	body.append_array(("\r\nContent-Disposition: form-data; name=\"%s\"; filename=\"%s\"" \
		% [fieldname, filepath.get_file()]).to_utf8_buffer())
	body.append_array(("\r\nContent-Type: %s\r\n\r\n" % mime_type).to_utf8_buffer())
	body.append_array(file_data)
	body.append_array("\r\n--".to_utf8_buffer())
	body.append_array(boundary)
	body.append_array("--\r\n".to_utf8_buffer())
	
	var request = HTTPRequest.new()
	
	var headers = [
	  "Content-Type: multipart/form-data; boundary=%s" % boundary.get_string_from_utf8()
  	]
	add_child(request)
	
	request.request_completed.connect(func(): request.queue_free())
	
	request.request_raw(url, headers, method, body)
	
	return request
	

func generate_boundary(filler_bytes: int):
	var boundary := "BOUNDARY".to_utf8_buffer()
	var filler = crypto.generate_random_bytes(filler_bytes)
	for byte in filler:
		boundary.append(filler_chars[byte%nfiller_chars].to_utf8_buffer()[0])
	return boundary
