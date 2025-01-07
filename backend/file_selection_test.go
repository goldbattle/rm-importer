package backend

import "testing"

func initFileSelection() FileSelection {
	folder1 := DocInfo{Id: "dir1", ParentId: "", IsFolder: true, Name: "folder1"}
	folder2 := DocInfo{Id: "dir2", ParentId: "", IsFolder: true, Name: "folder2"}
	folder3 := DocInfo{Id: "dir3", ParentId: "dir2", IsFolder: true, Name: "folder3"}
	file1 := DocInfo{Id: "f1", ParentId: "dir1", IsFolder: false, Name: "file1"}
	file2 := DocInfo{Id: "f2", ParentId: "dir2", IsFolder: false, Name: "file2"}
	file3 := DocInfo{Id: "f3", ParentId: "", IsFolder: false, Name: "file3"}
	c := make(map[DocId][]DocInfo)
	c[""] = []DocInfo{folder1, folder2, file3}
	c["dir1"] = []DocInfo{file1}
	c["dir2"] = []DocInfo{file2, folder3}
	return NewFileSelection(c)
}

func assertAllSelected(t *testing.T, f *FileSelection, folderId DocId) {
	for _, sel := range f.GetFolderSelection("") {
		if sel.Status == Selected {
			continue
		}
		t.Fatalf("After selecting root one of the items is not selected! %v", sel)
	}
}

func assertAllNotSelected(t *testing.T, f *FileSelection, folderId DocId) {
	for _, sel := range f.GetFolderSelection("") {
		if sel.Status == NotSelected {
			continue
		}
		t.Fatalf("After selecting root one of the items is not selected! %v", sel)
	}
}

func TestFolderSelectionRoot(t *testing.T) {
	f := initFileSelection()
	f.Select("", true)

	assertAllSelected(t, &f, "")
	assertAllSelected(t, &f, "dir1")
	assertAllSelected(t, &f, "dir2")

	f.Select("", false)

	assertAllNotSelected(t, &f, "")
	assertAllNotSelected(t, &f, "dir1")
	assertAllNotSelected(t, &f, "dir2")
}
