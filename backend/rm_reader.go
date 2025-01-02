package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DocId = string
type DocInfo struct {
	// Required
	Id       DocId
	ParentId DocId
	IsFolder bool
	Name     string

	// Optional
	Bookmarked   bool
	LastModified *time.Time
	FileType     *string
}

func parseDocsResponse(bytes []byte) ([]DocInfo, error) {
	data := []map[string]interface{}{}
	err := json.Unmarshal(bytes, &data)

	if err != nil {
		return nil, err
	}

	result := []DocInfo{}
	for _, item := range data {
		info := DocInfo{}

		// Required elements
		info.Id = item["ID"].(string)
		info.ParentId = item["Parent"].(string)
		info.IsFolder = bool(item["Type"].(string) == "CollectionType")
		info.Name = item["VissibleName"].(string)

		// Optional elements
		if _, ok := item["ModifiedClient"]; ok {
			last, err := time.Parse("2006-01-02T15:04:05.000Z", item["ModifiedClient"].(string))
			if err == nil {
				info.LastModified = &last
			}
		}

		if t, ok := item["FileType"].(string); ok {
			info.FileType = &t
		}

		info.Bookmarked = item["Bookmarked"].(bool)

		result = append(result, info)
	}

	return result, nil
}

func readDocs(tablet_addr string) ([]DocInfo, error) {
	if !IsIpValid(tablet_addr) {
		return nil, fmt.Errorf("readDocs error: the IP address is invalid")
	}

	directories := []string{""}
	result := []DocInfo{}

	for len(directories) > 0 {
		id := directories[0]
		directories = directories[1:]

		url := "http://" + tablet_addr + "/documents/" + id
		content_type := "application/json"

		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Post(url, content_type, &bytes.Buffer{})
		if err != nil {
			return nil, err
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		elements, err := parseDocsResponse(respBytes)
		if err != nil {
			return nil, err
		}

		for _, element := range elements {
			if element.IsFolder {
				directories = append(directories, element.Id)
			}
		}

		result = append(result, elements...)
	}

	return result, nil
}

type RmReader struct {
	/* For an collection with DocId 'id', map[id] stores elements in that folder.
	   For a root element the 'id' is empty
	*/
	items map[DocId][]DocInfo
}

/*
Reads all the items the rM tablet and stores them.
Past items are cleared in case read() was called previously.
*/
func (r *RmReader) Read(tablet_addr string) error {
	r.items = make(map[string][]DocInfo)
	docs, err := readDocs(tablet_addr)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		r.items[doc.ParentId] = append(r.items[doc.ParentId], doc)
	}

	return nil
}

func (r *RmReader) GetFolder(id DocId) []DocInfo {
	if items, ok := r.items[id]; ok {
		return items
	}
	return []DocInfo{}
}
