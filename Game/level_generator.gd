@icon("res://icons/hammericon.png")
class_name LevelGenerator extends Node

func generate_nodes(json_string: String) -> Node:
	var json = JSON.parse_string(json_string)
	if json is not Array:
		return null
	if !json.all(checkItem):
		return null
	var level = Node.new()
	for item in json:
		var region = Node2D.new()
		var byte_pool = Marshalls.base64_to_raw(item["regionImage"])
		var img = Image.new()
		img.load_png_from_buffer(byte_pool)
		var tex_rect = TextureRect.new()
		tex_rect.texture = ImageTexture.create_from_image(img)
		region.add_child(tex_rect)
		var collision = CollisionPolygon2D.new()
		var mesh = item["mesh"] as Array
		var vectormesh = mesh.map(func(v: Dictionary): return Vector2(v["x"], v["y"]))
		collision.polygon = vectormesh
		var col = StaticBody2D.new()
		col.add_child(collision)
		region.add_child(col)
		region.position = Vector2(item["cornerX"], item["cornerY"])
		level.add_child(region)
	return level

func checkItem(item: Variant) -> bool:
	return item is Dictionary and item.get("regionImage") is String and item.get("mesh") is Array \
		and item["mesh"].all(func(m): return m is Dictionary and m.has_all(["x", "y"])) and item.has_all(["cornerX", "cornerY"])
