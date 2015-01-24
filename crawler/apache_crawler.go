package crawler

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/PuerkitoBio/goquery"
)

type ApacheCrawler struct {
	BaseCrawler
}

func (crawler *ApacheCrawler) Crawl() error {
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

		headers := crawler.crawlHeaders()

		// Run through each row and extract data
		doc.Find("tr").Each(func(i int, s *goquery.Selection) {
			tds := s.Find("td")

			// Row is empty
			if text := strings.TrimSpace(tds.Text()); text == "" || text == "Parent Directory   -" {
				return
			}

			// Row has incorrect structure
			if tds.Size() != len(headers) {
				return
			}

			// Row has correct structure
			item := &elasticsearch.IndexItem{}
			tds.Each(func(i int, s *goquery.Selection) {
				// Ignore the img field
				if headers[i] == "img" {
					return
				}

				text := strings.TrimSpace(s.Text())
				if text == "" {
					return
				}

				if headers[i] == "Name" {
					href, _ := s.Find("a").First().Attr("href")
					item.Path, _ = url.QueryUnescape(href)
					item.Name = filepath.Base(item.Path)
				} else if headers[0] == "Size" {
					item.Size = ParseSize(text)

				} else if headers[i] == "LastModifiedAt" {
					date, err := ApacheParseDate(text)
					if err != nil {
						errs <- err
						return
					}
					item.LastModifiedAt = date
				}
			})
			itemsToIndex <- item
		})
		close(itemsToIndex)
		close(errs)
	}()

	return <-errs
}

func (crawler *ApacheCrawler) crawlHeaders() []string {
	var fields []string
	// Find table header and determine available fields
	crawler.Doc.Find("tr th").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())

		if s.Children().Is("img") {
			fields = append(fields, "img")
		} else if text != "" {
			switch text {
			case "Last modified":
				text = "LastModifiedAt"
			}
			fields = append(fields, text)
		}
	})
	return fields
}
