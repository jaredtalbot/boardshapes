class_name OwnedAfterImageEmitter extends AfterImageEmitter

func _ready():
	super()
	call_deferred("_link_with_owner_emitter")

func _link_with_owner_emitter():
	var ancestor: Node = get_parent()
	while(ancestor != null and not ancestor.is_in_group("EmitterOwners")):
		ancestor = ancestor.get_parent()
	
	if ancestor != null:
		var children = ancestor.get_children()
		for node in children:
			if node is AfterImageEmitter:
				node.emitted_after_image.connect(emit_after_image)
				return
