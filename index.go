package main

import (
	"context"
	"database/sql"
	"reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic"
	loremipsum "gopkg.in/loremipsum.v1"
)

var loremIpsumGenerator = loremipsum.New()

type Record struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Text string `json:"text"`
}

func InitSphinxConnection(port string) (*sql.DB, error) {
	return sql.Open("mysql", "root:@tcp(localhost:"+port+")/test")
}

func InitElasticConnection() (*elastic.Client, error) {
	return elastic.NewClient()
}

func InsertRandomRecordToSphinx(db *sql.DB, recordCount int) error {
	var err error
	var rows *sql.Rows
	var query string = "REPLACE INTO test_v1 (id, name, text) VALUES"
	var currentQuery string

	var i int
	for i = 0; i < recordCount; i++ {
		currentQuery = " ('" + strconv.Itoa(i+1) + "', '" + loremIpsumGenerator.Word() +
			"', '" + loremIpsumGenerator.Words(20) + "')"

		rows, err = db.Query(query + currentQuery)
		rows.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func InsertRandomRecordToElastic(client *elastic.Client, recordCount int) error {
	var err error
	var indexExists bool
	indexExists, err = client.IndexExists("test").Do(context.Background())
	if err != nil {
		return err
	}
	if !indexExists {
		_, err = client.CreateIndex("test").Do(context.Background())
		if err != nil {
			return err
		}
	}

	var bulk *elastic.BulkService
	bulk = client.
		Bulk().
		Index("test").
		Type("doc")

	var i int
	for i = 0; i < recordCount; i++ {
		var record = Record{
			Id:   strconv.Itoa(i + 1),
			Name: loremIpsumGenerator.Word(),
			Text: loremIpsumGenerator.Words(20),
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(record.Id).Doc(record))
	}
	_, err = bulk.Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func ReadSphinxIds(db *sql.DB, query string) ([]int, error) {
	var err error
	var rows *sql.Rows
	var id int
	var ids = make([]int, 0)

	if query != "" {
		query = "WHERE " + query
	}

	rows, err = db.Query("SELECT id FROM test_v1 " + query + " ORDER BY WEIGHT() DESC")
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func ReadElasticIds(client *elastic.Client) ([]string, error) {
	var err error
	var searchResult *elastic.SearchResult
	var ids = make([]string, 0)
	var termQuery *elastic.TermQuery
	var item *Record
	var tmp interface{}
	var ok bool
	termQuery = elastic.NewTermQuery("text", loremIpsumGenerator.Word())
	searchResult, err = client.Search().
		Index("test").
		Query(termQuery).
		From(0).
		Size(10).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	for _, tmp = range searchResult.Each(reflect.TypeOf(item)) {
		if item, ok = tmp.(*Record); ok {
			ids = append(ids, item.Id)
		}
	}
	return ids, nil
}

func main() {

}
