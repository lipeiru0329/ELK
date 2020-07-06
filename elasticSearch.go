package elk

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strconv"
)

var (
	subject   Subject
	indexName = "subject"
	servers   = []string{"http://localhost:9200/"}
)

type Subject struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
}

const mapping = `
{
    "mappings": {
        "properties": {
            "id": {
                "type": "long"
            },
            "title": {
                "type": "text"
            },
            "genres": {
                "type": "keyword"
            }
        }
    }
}`

func elasticSearch()  {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(servers...))
	if err != nil {
		panic(err)
	}

	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		_, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
	}

	subject = Subject{
		ID:     1,
		Title:  "肖恩克的救赎",
		Genres: []string{"犯罪", "剧情"},
	}

	// 写入
	doc, err := client.Index().
		Index(indexName).
		Id(strconv.Itoa(subject.ID)).
		BodyJson(subject).
		Refresh("wait_for").
		Do(ctx)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed with id=%v, type=%s\n", doc.Id, doc.Type)
	subject = Subject{
		ID:     2,
		Title:  "千与千寻",
		Genres: []string{"剧情", "喜剧", "爱情", "战争"},
	}
	fmt.Println(string(subject.ID))
	doc, err = client.Index().
		Index(indexName).
		Id(strconv.Itoa(subject.ID)).
		BodyJson(subject).
		Refresh("wait_for").
		Do(ctx)

	if err != nil {
		panic(err)
	}

}
