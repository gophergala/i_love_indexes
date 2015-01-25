package crawler

import (
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/PuerkitoBio/goquery"
	"github.com/Scalingo/go-workers"
)

type BaseCrawler struct {
	relativePath string
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
		item.Path = crawler.relativePath + "/" + strings.Trim(item.Path, "/")
		if item.Size == -1 {
			item.MIMEType = "directory"
			workers.Enqueue("index-crawler", "CrawlWorker", []string{crawler.IndexOf.Id, item.Path})
		}
		item.LastSeenAt = time.Now()
		item.Name, _ = url.QueryUnescape(filepath.Base(item.Path))
		item.SetEscapedName()
		item.URL = crawler.IndexOf.URL() + "/" + item.Path
		item.IndexOfId = crawler.IndexOf.Id
		elasticsearch.Index(item)
	}
}
