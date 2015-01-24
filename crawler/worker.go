package crawler

import (
	"log"

	"github.com/Scalingo/go-workers"
)

func CrawlWorker(msg *workers.Msg) {
	params, err := msg.Args().Array()
	if err != nil {
		log.Println("param is not a string:", err)
	}
	url := params[0].(string)
	// id := params[0].(string)
	Crawl(url)
}
