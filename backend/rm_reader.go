package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
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
	Path         *string
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

		if t, ok := item["fileType"].(string); ok {
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
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for len(directories) > 0 {
		id := directories[0]
		directories = directories[1:]

		url := "http://" + tablet_addr + "/documents/" + id

		req, err := http.NewRequest(http.MethodPost, url, &bytes.Buffer{})
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("User-Agent", "Mozilla/5.0 (U; Linux x86_64; en-US) Gecko/20100101 Firefox/133.0")

		resp, err := client.Do(req)
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
	children map[DocId][]DocInfo

	docById map[DocId]DocInfo
}

/*
Reads all the items the rM tablet and stores them.
Past items are cleared in case read() was called previously.
*/
func (r *RmReader) Read(tablet_addr string) error {
	r.children = make(map[DocId][]DocInfo)
	r.docById = make(map[DocId]DocInfo)

	docs, err := readDocs(tablet_addr)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		r.children[doc.ParentId] = append(r.children[doc.ParentId], doc)
	}

	for _, doc := range docs {
		r.docById[doc.Id] = doc
	}

	return nil
}

func (r *RmReader) GetFolder(id DocId) []DocInfo {
	if items, ok := r.children[id]; ok {
		return items
	}
	return []DocInfo{}
}

func (r *RmReader) GetElementsByIds(ids []DocId) []DocInfo {
	result := []DocInfo{}
	for _, id := range ids {
		result = append(result, r.docById[id])
	}
	return result
}

func (r *RmReader) GetElementById(id DocId) DocInfo {
	return r.docById[id]
}

func (r *RmReader) GetChildren() map[DocId][]DocInfo {
	return r.children
}

func (r *RmReader) GetPath(id DocId) string {
	l := []string{}
	for id != "" {
		item := r.docById[id]
		l = append(l, item.Name)
		id = item.ParentId
	}
	slices.Reverse(l)
	return strings.Join(l[:], "/")
}

func (r *RmReader) GetPaths(ids []DocId) []string {
	result := []string{}
	for _, id := range ids {
		result = append(result, r.GetPath(id))
	}
	return result
}
