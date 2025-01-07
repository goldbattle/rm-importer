package backend

type SelectionStatus = int

const (
	NotSelected   SelectionStatus = iota // default value
	Indeterminate SelectionStatus = iota
	Selected      SelectionStatus = iota
)

type SelectionInfo struct {
	Id     DocId
	Status SelectionStatus
}

type FileSelection struct {
	Children map[DocId][]DocId
	Parent   map[DocId]DocId

	DocSelection map[DocId]SelectionStatus
	Files        map[DocId]bool
	checkedFiles int
}

func NewFileSelection(c map[DocId][]DocInfo) FileSelection {
	m := make(map[DocId][]DocId)
	for parentId, items := range c {

		ids := []DocId{}
		for _, item := range items {
			ids = append(ids, item.Id)
		}

		m[parentId] = ids
	}

	parent := make(map[DocId]DocId)
	for parentId, items := range c {
		for _, item := range items {
			parent[item.Id] = parentId
		}
	}

	files := make(map[DocId]bool)
	for _, items := range c {
		for _, item := range items {
			if !item.IsFolder {
				files[item.Id] = true
			}
		}
	}

	return FileSelection{m, parent, make(map[string]SelectionStatus), files, 0}
}

func (f *FileSelection) setDocSelection(id DocId, selection bool) {
	var status SelectionStatus
	if !selection {
		status = NotSelected
	} else {
		status = Selected
	}

	if _, ok := f.Files[id]; ok {
		if f.DocSelection[id] == Selected && !selection {
			f.checkedFiles -= 1
		} else if f.DocSelection[id] == NotSelected && selection {
			f.checkedFiles += 1
		}
	}

	f.DocSelection[id] = status
}

func (f *FileSelection) dfs(item DocId, selection bool) {
	f.setDocSelection(item, selection)

	if ids, ok := f.Children[item]; ok {
		for _, child := range ids {
			f.dfs(child, selection)
		}
	}
}

func (f *FileSelection) updateParents(id DocId) {
	for {
		if ids, ok := f.Children[id]; ok {
			s, i := 0, 0
			for _, id := range ids {
				if f.DocSelection[id] == Selected {
					s += 1
				}
				if f.DocSelection[id] == Indeterminate {
					i += 1
				}
			}
			if s == 0 && i == 0 {
				f.DocSelection[id] = NotSelected
			} else if s == len(ids) {
				f.DocSelection[id] = Selected
			} else {
				f.DocSelection[id] = Indeterminate
			}
		}

		if id == "" {
			break
		}
		id = f.Parent[id]
	}
}

func (f *FileSelection) Select(id DocId, selection bool) {
	/* Set the selection value to all items in the subtree of 'item'. */
	f.dfs(id, selection)

	/* Number of selectedOrIndeterminate children could change,
	   which could in turn change the state of the parent itself.
	   This method iteratively updates parents' states up to root. */
	f.updateParents(id)
}

func (f *FileSelection) GetFolderSelection(id DocId) []SelectionInfo {
	result := []SelectionInfo{}
	for _, id := range f.Children[id] {
		result = append(result, SelectionInfo{id, f.DocSelection[id]})
	}
	return result
}

func (f *FileSelection) GetItemSelection(id DocId) SelectionInfo {
	return SelectionInfo{id, f.DocSelection[id]}
}

func (f *FileSelection) GetCheckedFiles() []DocId {
	result := []DocId{}
	for id := range f.Files {
		if f.DocSelection[id] == Selected {
			result = append(result, id)
		}
	}
	return result
}

func (f *FileSelection) GetCheckedFilesCount() int {
	return f.checkedFiles
}
