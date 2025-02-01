package backend

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

/*
Graph sketch

			 root
	/          |           \

dir1         dir2           file3
file1      dir3 file2
*/
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
	for _, sel := range f.GetFolderSelection(folderId) {
		if sel.Status == Selected {
			continue
		}
		t.Fatalf("After selecting root one of the items is not selected! %v", sel)
	}
}

func assertAllNotSelected(t *testing.T, f *FileSelection, folderId DocId) {
	for _, sel := range f.GetFolderSelection(folderId) {
		if sel.Status == NotSelected {
			continue
		}
		t.Fatalf("After selecting root one of the items is not selected! %v", sel)
	}
}

func assertOne(t *testing.T, f *FileSelection, id DocId, expected SelectionStatus) {
	result := f.GetItemSelection(id)
	if result.Status != expected {
		t.Fatalf("assertOne failed! id: %s, result.Status: %v, expected: %v", id, result.Status, expected)
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

func TestSelectionSub(t *testing.T) {
	f := initFileSelection()

	f.Select("", true)
	f.Select("dir1", false)
	f.Select("dir2", true)

	assertOne(t, &f, "", Indeterminate)
	assertAllNotSelected(t, &f, "dir1")
	assertAllSelected(t, &f, "dir2")
	assertOne(t, &f, "f3", Selected)
}

func TestParentUpdate(t *testing.T) {
	f := initFileSelection()
	f.Select("f1", true)
	f.Select("f2", true)
	f.Select("f3", true)
	assertOne(t, &f, "", Indeterminate)
	assertOne(t, &f, "dir1", Selected)
	assertOne(t, &f, "dir2", Indeterminate)
	assertOne(t, &f, "dir3", NotSelected)
}

func TestParentUpdate2(t *testing.T) {
	f := initFileSelection()
	f.Select("f1", true)
	f.Select("f2", true)
	f.Select("f3", true)
	f.Select("dir3", true)
	assertOne(t, &f, "", Selected)
	assertOne(t, &f, "dir1", Selected)
	assertOne(t, &f, "dir2", Selected)
	assertOne(t, &f, "dir3", Selected)
}

func TestSelectSimple(t *testing.T) {
	f := initFileSelection()
	f.Select("f1", true)
	assertOne(t, &f, "f1", Selected)
	assertOne(t, &f, "f2", NotSelected)
	assertOne(t, &f, "f3", NotSelected)
}

func TestCheckedFilesRoot(t *testing.T) {
	f := initFileSelection()
	f.Select("", true)

	result := f.GetCheckedItems()
	expected := []DocId{"f1", "f2", "f3"}
	slices.Sort(result)
	slices.Sort(expected)
	if !cmp.Equal(result, expected) {
		t.Fatalf("CheckedFilesRoot failed! %v", cmp.Diff(result, expected))
	}
}

func TestCheckedFilesSub(t *testing.T) {
	f := initFileSelection()
	f.Select("", false)
	f.Select("dir2", true)
	f.Select("f3", true)

	result := f.GetCheckedItems()
	expected := []DocId{"f2", "f3"}
	slices.Sort(result)
	slices.Sort(expected)
	if !cmp.Equal(result, expected) {
		t.Fatalf("CheckedFilesSub failed! %v", cmp.Diff(result, expected))
	}
}
