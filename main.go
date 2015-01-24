package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/GopherGala/i-love-indexes/api"
	"github.com/GopherGala/i-love-indexes/crawler"
	"github.com/Scalingo/go-workers"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/shaoshing/train"
)

var (
	testURL string
)

func main() {
	isWeb := flag.Bool("web", false, "run web server")
	isCrawler := flag.Bool("crawler", false, "run async crawler")
	isTester := flag.Bool("tester", false, "add a url in the queue")
	flag.StringVar(&testURL, "test-url", "", "url to test async crawling")
	flag.Parse()
	if *isWeb {
		mainWebServer()
	} else if *isCrawler {
		mainCrawler()
	} else if *isTester {
		mainTester()
	} else {
		log.Fatalln("Invalid type of process, precis -web, -crawler or -tester")
	}
}

func setupTrain(router *mux.Router) {
	train.SetFileServer()
	router.Handle(train.Config.AssetsUrl+"{any:.*}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		train.ServeRequest(w, r)
	}))
}

func mainWebServer() {
	router := mux.NewRouter()
	router.Handle("/api/{any:.*}", api.NewAPI())
	setupTrain(router)

	staticHandler := negroni.Classic()
	router.Handle("/{any:.*}", staticHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	log.Println("Listening on", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}

func mainCrawler() {
	configureWorkers()
	workers.Process("index-crawler", crawler.CrawlWorker, 10)
	workers.Run()
}

func mainTester() {
	configureWorkers()
	workers.Enqueue("index-crawler", "CrawlWorker", testURL)
}

func configureWorkers() {
	workers.Configure(map[string]string{
		"server":   redisHost(),
		"database": "0",
		"pool":     "30",
		"process":  "1",
	})
}

func redisHost() string {
	redisURL := "redis://localhost:6379"
	if os.Getenv("REDIS_URL") != "" {
		redisURL = os.Getenv("REDIS_URL")
	}
	u, err := url.Parse(redisURL)
	if err != nil {
		log.Fatalln(err)
	}
	return u.Host
}
