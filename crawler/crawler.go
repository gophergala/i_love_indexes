package crawler

import (
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/errgo.v1"
)

var (
	nginxServerRegexp      = regexp.MustCompile(`^nginx.*$`)
	apacheServerRegexp     = regexp.MustCompile(`^Apache.*$`)
	lighthttpdServerRegexp = regexp.MustCompile(`^lighttpd.*$`)
)

type Crawler interface {
	Crawl() error
}

type BaseCrawler struct {
	IndexOfId string
	Doc       *goquery.Document
}

func CrawlerFromUrl(url string, id string) (Crawler, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	server := res.Header.Get("Server")
	if nginxServerRegexp.MatchString(server) {
		return &NginxCrawler{BaseCrawler{id, doc}}, nil
	} else if lighthttpdServerRegexp.MatchString(server) {
		return &LighttpdCrawler{BaseCrawler{id, doc}}, nil
	} else if apacheServerRegexp.MatchString(server) {
		return &ApacheCrawler{BaseCrawler{id, doc}}, nil
	} else {
		return nil, errgo.Newf("Unknown 'Server' header: %v", server)
	}
}
