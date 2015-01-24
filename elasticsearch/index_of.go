package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	elastigo "github.com/mattbaird/elastigo/lib"
	"gopkg.in/errgo.v1"
)

var (
	AlreadyIndexedErr = errors.New("index of is already indexed")
)

type IndexOf struct {
	Id        string    `json:"_id,omitempty"`
	Host      string    `json:"host"`
	Scheme    string    `json:"scheme"`
	Path      string    `json:"path"`
	CrawledAt time.Time `json:"crawled_at"`
}

func (i *IndexOf) Type() string {
	return "index_of"
}

func (i *IndexOf) GetId() string {
	return i.Id
}

func (i *IndexOf) SetId(id string) {
	i.Id = id
}

func (i *IndexOf) URL() string {
	return i.Scheme + "://" + i.Host + i.Path
}

func FindIndexOf(id string) (*IndexOf, error) {
	res, err := defaultConn.Get(defaultIndex, "index_of", id, nil)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	var indexOf *IndexOf
	err = json.Unmarshal(*res.Source, &indexOf)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return indexOf, nil
}

func (i *IndexOf) Index() error {
	res, err := elastigo.Search(defaultIndex).Type(i.Type()).Query(elastigo.Query().Term("host", i.Host).Term("scheme", i.Scheme)).Result(defaultConn)
	if err != nil {
		if err == elastigo.RecordNotFound {
			return Index(i)
		}
		return errgo.Mask(err)
	}
	if res.Hits.Len() == 1 {
		return AlreadyIndexedErr
	}
	return Index(i)
}

type IndexItem struct {
	Id             string    `json:"_id,omitempty"`
	IndexOfId      string    `json:"index_of_id"`
	Name           string    `json:"name"`
	UpdatedAt      time.Time `json:"updated_at"`
	Size           int64     `json:"size,omitempty"`
	MIMEType       string    `json:"mime_type"`
	Description    string    `json:"description,omitempty"`
	Path           string    `json:"path"`
	LastSeenAt     time.Time `json:"last_seen_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
}

func (i *IndexItem) Type() string {
	return "index_item"
}

func (i *IndexItem) GetId() string {
	return i.Id
}

func (i *IndexItem) SetId(id string) {
	i.Id = id
}

func SearchIndexItemsPerName(name string) []*IndexItem {
	// Fuzzy search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"fuzzy_like_this_field": map[string]interface{}{
				"name": map[string]string{
					"like_text": name,
				},
			},
		},
	}

	var items []*IndexItem
	var item *IndexItem

	res, err := defaultConn.Search(defaultIndex, "index_item", nil, query)
	if err != nil {
		fmt.Println("fuzzy search err:", err)
	}

	for _, h := range res.Hits.Hits {
		json.Unmarshal(*h.Source, &item)
		items = append(items, item)
	}

	return items
}
