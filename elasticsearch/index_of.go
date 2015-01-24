package elasticsearch

import "time"

type IndexOf struct {
	URL       string    `json:"url"`
	CrawledAt time.Time `json:"crawled_at"`
}

type IndexItem struct {
	Name        string    `json:"name"`
	UpdatedAt   time.Time `json:"updated_at"`
	Size        int64     `json:"size,omitempty"`
	Description string    `json:"description,omitempty"`
	Path        string    `json:"path"`
	LastSeenAt  time.Time `json:"last_seen_at"`
}
