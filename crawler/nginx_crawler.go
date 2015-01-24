package crawler

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
)

type NginxCrawler struct {
	BaseCrawler
}

var (
	spaceRegexp = regexp.MustCompile(`\s\s+`)
)

func (crawler *NginxCrawler) Crawl() error {
	errs := make(chan error)
	itemsToIndex := make(chan *elasticsearch.IndexItem)
	go func() {
		for item := range itemsToIndex {
			item.IndexOfId = crawler.IndexOfId
			err := elasticsearch.Index(item)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(item.Id)
		}
	}()

	go func() {
		doc := crawler.Doc

		// Run through each row and extract data
		pre := doc.Find("pre").First()
		as := pre.Find("a").Nodes
		entries := strings.Split(pre.Text(), "\n")
		var err error
		for i, entry := range entries {
			entry = strings.TrimSpace(entry)
			if len(entry) == 0 {
				continue
			}
			item := &elasticsearch.IndexItem{}
			item.Path, _ = url.QueryUnescape(as[i].Attr[0].Val)
			item.Name = filepath.Base(item.Path)
			if strings.Contains(entry, "../") {
				continue
			}
			fields := strings.Split(spaceRegexp.ReplaceAllString(entry[51:], "\t"), "\t")
			item.LastModifiedAt, err = ApacheParseDate(fields[0])
			if err != nil {
				errs <- err
			}
			item.Size = mustInt64(fields[1])
			itemsToIndex <- item
		}

		close(itemsToIndex)
		close(errs)
	}()
	return <-errs
}
