package backend

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type RmExport struct {
	Format             string // 'pdf' or 'rmdoc'
	Location           string // path to the folder to export
	WrappingFolderName string
}

/* Returns 0 if all items were succesfully exported, otherwise returns the index of an item that errored first. */
/* 'export_from' is the index in items list of the first item to be exported.
   All items with index less than export_from are ignored. */
func (r *RmExport) ExportMultiple(ctx context.Context, tablet_addr string, items []DocInfo, export_from int) int {
	client := http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
		},
		Timeout: 5 * time.Minute,
	}

	runtime.LogInfof(ctx, "Export format: %v", r.Format)
	runtime.LogInfof(ctx, "In export location, using a wrapper folder with a name: %v", r.WrappingFolderName)

	// possible states of export: downloading, finished, error
	for i, item := range items {
		if i < export_from {
			continue
		}
		runtime.EventsEmit(ctx, "downloading", item.Id)
		runtime.LogInfof(ctx, "Downloading file with id=%v", item.Id)

		err := r.export(&client, tablet_addr, item)

		if err == nil {
			runtime.LogInfof(ctx, "Finished downloading file with id=%v", item.Id)
			runtime.EventsEmit(ctx, "finished", item.Id)
		} else {
			runtime.LogInfof(ctx, "Error while downloading file with id=%v, error=%v", item.Id, err.Error())
			runtime.EventsEmit(ctx, "error", item.Id, err.Error())
			return i
		}
	}

	return 0
}

func (r *RmExport) export(client *http.Client, tablet_addr string, item DocInfo) error {
	if item.IsFolder {
		return nil
	}

	out, err := r.createFile(r.WrappingFolderName, item)
	if err != nil {
		return err
	}
	defer out.Close()

	url := "http://" + tablet_addr + "/download/" + item.Id + "/" + r.Format
	resp, err := client.Get(url)

	if resp.StatusCode != 200 {
		return fmt.Errorf("tablet returned HTTP code %d", resp.StatusCode)
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

func (r *RmExport) createFile(folderName string, item DocInfo) (*os.File, error) {
	itemPath := *item.Path
	if path.Ext(itemPath) != "."+r.Format {
		itemPath = itemPath + "." + r.Format
	}

	path := filepath.Join(filepath.FromSlash(r.Location), folderName, itemPath)
	dir, _ := filepath.Split(path)

	/* Permission 0755: The owner can read, write, execute.
	   Everyone else can read and execute but not modify the file.*/
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return out, nil
}
