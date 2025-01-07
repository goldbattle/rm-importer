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
	Children     map[DocId][]DocId
	DocSelection map[DocId]SelectionStatus
	Parent       map[DocId]DocId
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

	return FileSelection{m, make(map[string]SelectionStatus), parent}
}

func (f *FileSelection) setDocSelection(id DocId, selection bool) {
	var status SelectionStatus
	if !selection {
		status = NotSelected
	} else {
		status = Selected
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

func (f *FileSelection) updateParents(id DocId, isFolder bool) {
	for {
		if isFolder {
			ids := f.Children[id]
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
		isFolder = true
	}
}

func (f *FileSelection) Select(item DocInfo, selection bool) {
	/* Set the selection value to all items in the subtree of 'item'. */
	f.dfs(item.Id, selection)

	/* Number of selectedOrIndeterminate children could change,
	   which could in turn change the state of the parent itself.
	   This method iteratively updates parents' states up to root. */
	f.updateParents(item.Id, item.IsFolder)
}

func (f *FileSelection) GetFolderSelection(id DocId) []SelectionInfo {
	result := []SelectionInfo{}
	for _, id := range f.Children[id] {
		result = append(result, SelectionInfo{id, f.DocSelection[id]})
	}
	return result
}
