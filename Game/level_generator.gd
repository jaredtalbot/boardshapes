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
		var sprite = Sprite2D.new()
		sprite.centered = false
		sprite.texture = ImageTexture.create_from_image(img)
		region.add_child(sprite)
		var collision = CollisionPolygon2D.new()
		var mesh = item["mesh"] as Array
		var vectormesh = mesh.map(func(v: Dictionary): return Vector2(v["x"], v["y"]))
		collision.polygon = vectormesh
		var col = StaticBody2D.new()
		col.add_child(collision)
		region.add_child(col)
		region.position = Vector2(item["cornerX"], item["cornerY"])
		var color = item["regionColorString"]
		match color:
			"Red":
				col.add_to_group("Red")
				if ProjectSettings.get_setting("rendering/environment/defaults/color_blind_mode") == true:
					sprite.material = ShaderMaterial.new()
					sprite.material.shader = load("res://colorblind_filter.gdshader")
					sprite.material.set("shader_parameter/tile_size", 2)
					sprite.material.set("shader_parameter/pattern", load("res://red_cb.png"))
			"Green":
				col.add_to_group("Green")
				if ProjectSettings.get_setting("rendering/environment/defaults/color_blind_mode") == true:
					sprite.material = ShaderMaterial.new()
					sprite.material.shader = load("res://colorblind_filter.gdshader")
					sprite.material.set("shader_parameter/tile_size", 2)
					sprite.material.set("shader_parameter/pattern", load("res://green_cb.png"))
			"Blue":
				col.add_to_group("Blue")
				if ProjectSettings.get_setting("rendering/environment/defaults/color_blind_mode") == true:
					sprite.material = ShaderMaterial.new()
					sprite.material.shader = load("res://colorblind_filter.gaadshader")
					sprite.material.set("shader_parameter/tile_size", 2)
					sprite.material.set("shader_parameter/pattern", load("res://blue_cb.png"))
			"Black":
				if RenderingServer.get_default_clear_color() == Color(0, 0, 0, 1):
					sprite.material = ShaderMaterial.new()
					sprite.material.shader = load("res://color_invert.gdshader")
		level.add_child(region)
	return level

func checkItem(item: Variant) -> bool:
	return item is Dictionary and item.get("regionImage") is String and item.get("mesh") is Array \
		and item["mesh"].all(func(m): return m is Dictionary and m.has_all(["x", "y"])) and item.has_all(["cornerX", "cornerY"])
