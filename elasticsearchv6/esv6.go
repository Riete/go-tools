package elasticsearchv6

import (
	"github.com/elastic/go-elasticsearch/v6"
	"log"
	"strings"
)

type ESv6 struct {
	Addresses []string
	Username  string
	Password  string
	Client    *elasticsearch.Client
}

func (es *ESv6) NewClient() {
	config := elasticsearch.Config{Addresses: es.Addresses}
	if es.Username != "" && es.Password != "" {
		config = elasticsearch.Config{
			Addresses: es.Addresses,
			Password:  es.Password,
			Username:  es.Username,
		}
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
	es.Client = client
}

func (es *ESv6) Search(index, queryDSL string) string {
	resp, err := es.Client.Search(
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithPretty(),
		es.Client.Search.WithBody(strings.NewReader(queryDSL)),
	)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	result := strings.Trim(resp.String(), "[200 OK] ")
	return result
}
