package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type DateTime struct {
	time.Time
}

func (t *DateTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000Z"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return nil
}

type IsFolderBool bool

func (f *IsFolderBool) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	if s == `"CollectionType"` {
		*f = false
		return nil
	} else if s == `"DocumentType"` {
		*f = true
		return nil
	}
	return fmt.Errorf("IsFolderBool unmarshal error: invalid API response %s", s)
}

type DocInfo struct {
	// Required
	Id       string
	ParentId string
	IsFolder bool
	Name     string

	// Optional
	Bookmarked   bool
	LastModified *time.Time
	FileType     *string
}

var tablet_addr = "10.11.99.1"

func ParseFilesResponse(bytes []byte) ([]DocInfo, error) {
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

func ReadFiles() []DocInfo {
	directories := []string{""}
	result := []DocInfo{}

	for len(directories) > 0 {
		dir := directories[0]
		directories = directories[1:]

		url := "http://" + tablet_addr + "/documents" + dir + "/"
		content_type := "application/json"

		resp, err := http.Post(url, content_type, &bytes.Buffer{})
		if err != nil {
			log.Fatal(err)
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		elements, err := ParseFilesResponse(respBytes)

		if err != nil {
			log.Fatal(err)
		}

		result = append(result, elements...)
	}

	return result
}
