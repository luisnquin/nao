package store

func (d *Box) GetHistoryOf(key string) ([]Change, error) {
	for k, v := range d.box.NaoSet {
		if k == key {
			return v.History, nil
		}
	}

	return nil, ErrNoteNotFound
}

func (d *Box) CleanHistoryOf(key string) error {
	for k, v := range d.box.NaoSet {
		if k == key {
			v.History = nil
		}

		d.box.NaoSet[k] = v

		return d.updateFile()
	}

	return ErrNoteNotFound
}

func (d *Box) CleanHistoryOfAll() error {
	for k, v := range d.box.NaoSet {
		v.History = nil

		d.box.NaoSet[k] = v
	}

	return d.updateFile()
}
