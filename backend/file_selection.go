package backend

type SelectionStatus = int

const (
	NotSelected   SelectionStatus = iota // default value
	Indeterminate SelectionStatus = iota
	Selected      SelectionStatus = iota
)

type children struct {
	ids                     []DocId
	selectedOrIndeterminate int
}

type FileSelection struct {
	children     map[DocId]children
	docSelection map[DocId]SelectionStatus
}

func New(c map[DocId][]DocInfo) FileSelection {
	m := make(map[DocId]children)
	for parentId, items := range c {

		ids := []DocId{}
		for _, item := range items {
			ids = append(ids, item.Id)
		}

		m[parentId] = children{ids, 0}
	}

	return FileSelection{m, make(map[string]SelectionStatus)}
}

func (f *FileSelection) setDocSelection(parentId DocId, id DocId, selection bool) {
	prev_status := f.docSelection[id]

	var status SelectionStatus
	if !selection {
		status = NotSelected
	} else {
		status = Selected
	}

	if id == "" {
		return
	}

	diff := 0
	if prev_status == NotSelected && status == Selected {
		diff = 1
	} else if (prev_status == Selected || prev_status == Indeterminate) && status == NotSelected {
		diff = -1
	}
	if diff != 0 {
		entry := f.children[parentId]
		entry.selectedOrIndeterminate += diff
		f.children[parentId] = entry
	}
}

func (f *FileSelection) dfs(parent DocId, item DocId, selection bool) {
	f.setDocSelection(parent, item, selection)

	if c, ok := f.children[item]; ok {
		for _, child := range c.ids {
			f.dfs(item, child, selection)
		}
	}
}

func (f *FileSelection) OnSelect(item DocInfo, selection bool) {
	// TODO what if item is root with empty id?
	/* Set the selection value to all items in the subtree of 'item'. */
	f.dfs(item.ParentId, item.Id, selection)

}
