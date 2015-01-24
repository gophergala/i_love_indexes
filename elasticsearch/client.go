package elasticsearch

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	elastigo "github.com/mattbaird/elastigo/lib"
)

var (
	defaultConn *elastigo.Conn
)

func init() {
	initDefaultConn()
	initIndexes()
}

func initDefaultConn() {
	u := os.Getenv("ELASTICSEARCH_URL")
	if u != "" {
		elasticsearchURL, err := url.Parse(u)
		if err != nil {
			log.Fatalln("Invlid URL:", u)
		}
		splittedHost := strings.Split(elasticsearchURL.Host, ":")
		defaultConn = &elastigo.Conn{
			Protocol:       elastigo.DefaultProtocol,
			Domain:         splittedHost[0],
			ClusterDomains: []string{splittedHost[0]},
			Port:           splittedHost[1],
			DecayDuration:  time.Duration(elastigo.DefaultDecayDuration * time.Second),
		}
	}
}

func initIndexes() {
	res, err := defaultConn.CreateIndex("i<3indexes")
	if err != nil {
		log.Fatalln(err)
	}
	if res.Exists {
		log.Println("Index i<3indexes already exists")
	}
}
