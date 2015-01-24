package elasticsearch

import "time"

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
