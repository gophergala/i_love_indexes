package elasticsearch

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type IndexItem struct {
	Id             string    `json:"_id,omitempty"`
	IndexOfId      string    `json:"index_of_id"`
	Name           string    `json:"name"`
	EscapedName    string    `json:"escaped_name"`
	Path           string    `json:"path"`
	UpdatedAt      time.Time `json:"updated_at"`
	Size           int64     `json:"size,omitempty"`
	URL            string    `json:"url",omitempty`
	MIMEType       string    `json:"mime_type,omitempty"`
	Description    string    `json:"description,omitempty"`
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

func (i *IndexItem) SetEscapedName() {
	r := strings.NewReplacer("-", " ", "_", " ", ".", " ")
	i.EscapedName = r.Replace(i.Name)
}

func SearchIndexItemsPerName(name string) []*IndexItem {
	// Fuzzy search query
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"fuzzy_like_this_field": map[string]interface{}{
	// 			"escaped_name": map[string]string{
	// 				"like_text": name,
	// 			},
	// 		},
	// 	},
	// }

	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"fuzzy": map[string]interface{}{
	// 			"escaped_name": name,
	// 		},
	// 	},
	// }

	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"fuzzy": map[string]interface{}{
	// 			"escaped_name": map[string]interface{}{
	// 				"value":     name,
	// 				"fuzziness": 2,
	// 			},
	// 		},
	// 	},
	// }

	// Match with fuzziness
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match": map[string]interface{}{
	// 			"escaped_name": map[string]interface{}{
	// 				"query":     name,
	// 				"fuzziness": "AUTO",
	// 			},
	// 		},
	// 	},
	// }

	// Full-text query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"escaped_name": map[string]interface{}{
					"query":     name,
					"fuzziness": 2,
					"type":      "phrase",
				},
			},
		},
	}

	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"bool": map[string]interface{}{
	// 			"should": []interface{}{
	// 				map[string]interface{}{
	// 					"match": map[string]interface{}{
	// 						"escaped_name": map[string]interface{}{
	// 							"query":     name,
	// 							"fuzziness": 0.5,
	// 						},
	// 					},
	// 				},
	// 				map[string]interface{}{
	// 					"match": map[string]interface{}{
	// 						"escaped_name": map[string]interface{}{
	// 							"query":     name,
	// 							"fuzziness": 2,
	// 							"type":      "phrase",
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"sort": []interface{}{
	// 		"_score",
	// 	},
	// }

	items := []*IndexItem{}
	var item *IndexItem

	res, err := defaultConn.Search(defaultIndex, "index_item", nil, query)
	if err != nil {
		fmt.Println("fuzzy search err:", err)
	}

	fmt.Printf("%+v\n", res.Hits)

	for _, h := range res.Hits.Hits {
		item = &IndexItem{}

		json.Unmarshal(*h.Source, &item)
		items = append(items, item)
	}

	return items
}

func (i *IndexItem) Index() error {
	i.UpdatedAt = time.Now()
	return Index(i)
}
