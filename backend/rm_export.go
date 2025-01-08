package backend

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RmExport struct {
	Format   string // 'pdf' or 'rmdoc'
	Location string // path to the folder to export
}

/*
Downloads all urls passed as arguments and passes them directly to 1Password CLI
*/

func (r *RmExport) Export(tablet_addr string, item DocInfo) error {
	url := "http://" + tablet_addr + "/download/" + item.Id + "/placeholder"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("op", "document", "create", "--title", item.Name, "--file-name", item.Name+".pdf", "--format", "json")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, resp.Body)
	}()

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()
	if err != nil {
		return "", errors.New(errb.String())
	}

	return outb.String(), nil
}

func (r *RmExport) ExportPdf(tablet_addr string, item DocInfo) error {
	url := "http://" + tablet_addr + "/download/" + item.Id + "/placeholder"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	dirPathSplit := []string{r.Location}
	dirPathSplit = append(dirPathSplit, strings.Split(*item.Path, "/")...)
	dirPath := filepath.Join(dirPathSplit...)

	/* permission 0755: The owner can read, write, execute. Everyone else can read and execute but not modify the file.*/
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	/*filePath := append(dirPathSplit, item.


	cmd := exec.Command("op", "document", "create", "--title", item.Name, "--file-name", item.Name+".pdf", "--format", "json")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, resp.Body)
	}()*/
}
