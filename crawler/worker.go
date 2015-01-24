package crawler

import (
	"log"

	"github.com/Scalingo/go-workers"
	"gopkg.in/errgo.v1"
)

func CrawlWorker(msg *workers.Msg) {
	params, err := msg.Args().Array()
	if err != nil {
		log.Println("param is not a string:", err)
		return
	}
	url := params[0].(string)
	// id := params[0].(string)
	crawler, err := CrawlerFromUrl(url)
	if err != nil {
		log.Println("failed to get crawler", errgo.Details(err))
		return
	}
	err = crawler.Crawl()
	if err != nil {
		log.Println("failed to crawl", errgo.Details(err))
		return
	}
}
