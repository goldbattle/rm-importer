package backend

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os/exec"
)

/*
Downloads all urls passed as arguments and passes them directly to 1Password CLI
*/

func ExportPdf(tablet_addr string, item DocInfo) (string, error) {
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
