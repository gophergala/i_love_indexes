package crawler

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Crawl(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("goquery:", err)
	}

	var fields []string

	// Find table header and determine available fields
	doc.Find("tr th").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())

		if s.Children().Is("img") {
			fields = append(fields, "img")
		} else if text != "" {
			fields = append(fields, text)
		}
	})

	var items []map[string]string

	// Run through each row and extract data
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")

		// Row is empty
		if text := strings.TrimSpace(tds.Text()); text == "" {
			return
		}

		// Row has incorrect structure
		if tds.Size() != len(fields) {
			return
		}

		// Row has correct structure
		data := map[string]string{}
		tds.Each(func(i int, s *goquery.Selection) {
			// Ignore the img field
			if fields[i] == "img" {
				return
			}

			data[fields[i]] = s.Text()
		})
		items = append(items, data)

		fmt.Println(data)
		fmt.Println("=============")
	})
}
