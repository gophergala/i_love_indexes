package elasticsearch

import (
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
	URL       string    `json:"url"`
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

func (i *IndexOf) Index() error {
	res, err := elastigo.Search(defaultIndex).Type(i.Type()).Query(elastigo.Query().Term("url", i.URL)).Result(defaultConn)
	if err != nil {
		if err == elastigo.RecordNotFound {
			return Index(i)
		}
		return errgo.Mask(err)
	}
	fmt.Println(res.Hits)
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
