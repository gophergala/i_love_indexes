package crawler

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/PuerkitoBio/goquery"
)

type LighttpdCrawler struct {
	BaseCrawler
}

func (crawler *LighttpdCrawler) Crawl() error {
	errs := make(chan error)
	itemsToIndex := make(chan *elasticsearch.IndexItem)
	go func() {
		for item := range itemsToIndex {
			item.IndexOfId = crawler.IndexOfId
			elasticsearch.Index(item)
		}
	}()

	go func() {
		doc := crawler.Doc

		// Run through each row and extract data
		doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
			tds := s.Find("td")

			item := &elasticsearch.IndexItem{}
			tds.Each(func(i int, s *goquery.Selection) {
				text := strings.TrimSpace(s.Text())
				if text == "" {
					return
				}

				class, _ := s.Attr("class")
				var err error
				switch class {
				case "n":
					href, _ := s.Find("a").First().Attr("href")
					item.Path, _ = url.QueryUnescape(href)
					item.Name = filepath.Base(item.Path)
				case "m":
					item.LastModifiedAt, err = LighttpdParseDate(text)
					if err != nil {
						errs <- err
					}
				case "s":
					item.Size = ParseSize(text)
				case "t":
					item.MIMEType = text
				}
			})
			itemsToIndex <- item
		})
		close(itemsToIndex)
		close(errs)
	}()

	return <-errs
}