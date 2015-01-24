package crawler

import (
	"net/url"
	"path/filepath"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/PuerkitoBio/goquery"
)

type BaseCrawler struct {
	itemsToIndex chan *elasticsearch.IndexItem
	IndexOf      *elasticsearch.IndexOf
	Doc          *goquery.Document
}

func (crawler *BaseCrawler) Start() {
	crawler.indexResults()
}

func (crawler *BaseCrawler) End() {
	close(crawler.itemsToIndex)
}

func (crawler *BaseCrawler) indexResults() {
	for item := range crawler.itemsToIndex {
		item.Name, _ = url.QueryUnescape(filepath.Base(item.Path))
		item.SetEscapedName()
		item.URL = crawler.IndexOf.URL() + "/" + item.Path
		item.IndexOfId = crawler.IndexOf.Id
		elasticsearch.Index(item)
	}
}
