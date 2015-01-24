package crawler

import (
	"log"

	"github.com/Scalingo/go-workers"
)

func CrawlWorker(msg *workers.Msg) {
	url, err := msg.Args().String()
	if err != nil {
		log.Println("param is not a string:", err)
	}
	Crawl(url)
}
