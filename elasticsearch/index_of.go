package elasticsearch

import (
	"encoding/json"
	"errors"
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
	indexOf.Id = id
	return indexOf, nil
}

func (i *IndexOf) Index() error {
	res, err := elastigo.Search(defaultIndex).Type(i.Type()).Query(elastigo.Query().Term("host", i.Host)).Result(defaultConn)
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
