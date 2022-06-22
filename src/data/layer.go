package data

func (d *Data) GetSet(key string) (string, error)

func (d *Data) ModifySet(key string, content string) error

func (d *Data) SearchInSets(pattern string) (string, Set, error)

func (d *Data) ListSets() map[string]Set

func (d *Data) ListSetWithHiddenContent() []Window

func (d *Data) GetMainNote() string

func (d *Data) ModifyMainNote(content string) error
