package data

func (d *Box) GetGroups() []string {
	return d.box.Groups
}

func (d *Box) NewGroup(name string) error {
	if d.GroupExists(name) {
		return ErrGroupAlreadyExists
	}

	d.box.Groups = append(d.box.Groups, name)

	return d.updateBoxFile()
}

func (d *Box) DeleteGroupWithRelated(name string) error {
	for i, group := range d.box.Groups {
		if group == name {
			d.box.Groups = append(d.box.Groups[:i], d.box.Groups[i+1:]...)

			for k, note := range d.box.NaoSet {
				if note.Group == name {
					delete(d.box.NaoSet, k)
				}
			}

			return d.updateBoxFile()
		}
	}

	return ErrGroupNotFound
}

func (d *Box) GroupExists(name string) bool {
	for _, g := range d.box.Groups {
		if g == name {
			return true
		}
	}

	return false
}

func (d *Box) ModifyGroupName(oldName, newName string) error {
	for i, group := range d.box.Groups {
		if group == oldName {
			d.box.Groups[i] = newName

			for k, note := range d.box.NaoSet {
				if note.Group == oldName {
					note.Group = newName

					d.box.NaoSet[k] = note
				}
			}

			return d.updateBoxFile()
		}
	}

	return ErrGroupNotFound
}

func (d *Box) DeleteGroup(name string) error {
	for i, group := range d.box.Groups {
		if group == name {
			d.box.Groups = append(d.box.Groups[:i], d.box.Groups[i+1:]...)

			for k, note := range d.box.NaoSet {
				if note.Group == name {
					note.Group = ""
					d.box.NaoSet[k] = note
				}
			}

			return d.updateBoxFile()
		}
	}

	return ErrGroupNotFound
}

func (d *Box) ModifyAssignedGroup(key, name string) error {
	if !d.GroupExists(name) {
		return ErrGroupNotFound
	}

	for k, v := range d.box.NaoSet {
		if k == key {
			v.Group = name
			d.box.NaoSet[k] = v

			return d.updateBoxFile()
		}
	}

	return ErrNoteNotFound
}

func (d *Box) RemoveFromAssignedGroup(key string) error {
	for k, v := range d.box.NaoSet {
		if k == key {
			v.Group = ""
			d.box.NaoSet[k] = v

			return d.updateBoxFile()
		}
	}

	return ErrNoteNotFound
}
