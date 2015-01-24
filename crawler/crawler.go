package crawler

import (
	"net/http"
	"regexp"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
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

func NewCrawler(indexOf *elasticsearch.IndexOf) (Crawler, error) {
	res, err := http.Get(indexOf.URL())
	if err != nil {
		return nil, errgo.Mask(err)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	server := res.Header.Get("Server")
	baseCrawler := BaseCrawler{
		itemsToIndex: make(chan *elasticsearch.IndexItem, 10),
		IndexOf:      indexOf,
		Doc:          doc,
	}
	go baseCrawler.Start()

	if nginxServerRegexp.MatchString(server) {
		return &NginxCrawler{baseCrawler}, nil
	} else if lighthttpdServerRegexp.MatchString(server) {
		return &LighttpdCrawler{baseCrawler}, nil
	} else if apacheServerRegexp.MatchString(server) {
		return &ApacheCrawler{baseCrawler}, nil
	} else {
		return nil, errgo.Newf("Unknown 'Server' header: %v", server)
	}
}
