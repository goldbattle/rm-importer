package backend

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	names := []string{"hello: 1.pdf", "project/v2", "NUL", "NUL-1", "..pdf", "allgood", "cool.epub"}
	expected := []string{"hello- 1.pdf", "project-v2", "NUL-1", "NUL-1", "..pdf", "allgood", "cool.epub"}

	for i, name := range names {
		if normalize(name) != expected[i] {
			t.Fatalf("normalize!=expected, name=%v, normalize=%v, expected=%v", name, normalize(name), expected[i])
		}
	}
}

type args struct {
	location   string
	folderName string
	itemPath   []string
	ext        string
}

func TestGetFilePath(t *testing.T) {
	argsList := []args{
		{
			location:   "/Desktop/folder1/",
			folderName: "myFolder: 1",
			itemPath:   []string{"file1"},
			ext:        "pdf",
		},
		{
			location:   "/Desktop/folder1/",
			folderName: "folder2",
			itemPath:   []string{"file:2.epub"},
			ext:        "rmdoc",
		},
		{
			location:   "/Desktop/folder1/",
			folderName: "",
			itemPath:   []string{"file"},
			ext:        "epub",
		},
	}

	expected := []string{
		"/Desktop/folder1/myFolder- 1/file1.pdf",
		"/Desktop/folder1/folder2/file-2.epub.rmdoc",
		"/Desktop/folder1/file.epub",
	}

	for i, a := range argsList {
		res, err := getFilePath(a.location, a.folderName, a.itemPath, a.ext)
		if err != nil {
			t.Fatal(err.Error())
		}

		if res != expected[i] {
			t.Fatalf("getFilePath: args=%v, res=%v, expected=%v", a, res, expected[i])
		}
	}
}

func TestGetFilePathUnique(t *testing.T) {
	paths := initPaths()
	argsList := []args{
		{
			location:   "/Desktop/folder1/",
			folderName: "myFolder: 1",
			itemPath:   []string{"file1"},
			ext:        "pdf",
		},
		{
			location:   "/Desktop/folder1/",
			folderName: "myFolder: 1",
			itemPath:   []string{"file1"},
			ext:        "rmdoc",
		},
		{
			location:   "/Desktop/folder1/",
			folderName: "myFolder: 1",
			itemPath:   []string{"file1"},
			ext:        "pdf",
		},
		{
			location:   "/another",
			folderName: "f",
			itemPath:   []string{"file1"},
			ext:        "pdf",
		},
		{
			location:   "/another",
			folderName: "f",
			itemPath:   []string{"file2"},
			ext:        "pdf",
		},
		{
			location:   "/another",
			folderName: "f",
			itemPath:   []string{"file2-1"},
			ext:        "pdf",
		},
		{
			location:   "/another",
			folderName: "f",
			itemPath:   []string{"file2"},
			ext:        "pdf",
		},
		{
			location:   "/loc",
			folderName: "",
			itemPath:   []string{"doc", "file1"},
			ext:        "rmdoc",
		},
		{
			location:   "/loc",
			folderName: "",
			itemPath:   []string{"doc", "file1"},
			ext:        "rmdoc",
		},
		{
			location:   "/loc",
			folderName: "",
			itemPath:   []string{"doc", "file1"},
			ext:        "rmdoc",
		},
	}

	expected := []string{
		"/Desktop/folder1/myFolder- 1/file1.pdf",
		"/Desktop/folder1/myFolder- 1/file1.rmdoc",
		"/Desktop/folder1/myFolder- 1/file1-1.pdf",
		"/another/f/file1.pdf",
		"/another/f/file2.pdf",
		"/another/f/file2-1.pdf",
		"/another/f/file2-1-1.pdf",
		"/loc/doc/file1.rmdoc",
		"/loc/doc/file1-1.rmdoc",
		"/loc/doc/file1-2.rmdoc",
	}

	for i, a := range argsList {
		res, err := paths.getFilePathUnique(a.location, a.folderName, a.itemPath, a.ext)
		if err != nil {
			t.Fatal(err.Error())
		}

		if res != expected[i] {
			t.Fatalf("getFilePathUnique: args=%v, res=%v, expected=%v", a, res, expected[i])
		}
	}
}
