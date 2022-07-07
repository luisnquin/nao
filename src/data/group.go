package data

func (d *Box) GetGroups() []string {
	return d.box.Groups
}

func (d *Box) NewGroup(name string) error {
	if d.GroupExists(name) {
		return ErrGroupAlreadyExists
	}

	d.box.Groups = append(d.box.Groups, name)

	return d.updateFile()
}

func (d *Box) DeleteGroupWithRelated(name string) error {
	for i, group := range d.box.Groups {
		if group == name {
			d.box.Groups = append(d.box.Groups[:i], d.box.Groups[i+1:]...)

			for k, set := range d.box.NaoSet {
				if set.Group == name {
					delete(d.box.NaoSet, k)
				}
			}

			return d.updateFile()
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

			for k, set := range d.box.NaoSet {
				if set.Group == oldName {
					set.Group = newName

					d.box.NaoSet[k] = set
				}
			}

			return d.updateFile()
		}
	}

	return ErrGroupNotFound
}

func (d *Box) DeleteGroup(name string) error {
	for i, group := range d.box.Groups {
		if group == name {
			d.box.Groups = append(d.box.Groups[:i], d.box.Groups[i+1:]...)

			for k, set := range d.box.NaoSet {
				if set.Group == name {
					set.Group = ""
					d.box.NaoSet[k] = set
				}
			}

			return d.updateFile()
		}
	}

	return ErrGroupNotFound
}
