package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/olivere/elastic"
)

var recordCount = 10000
var searchCount = 1000

func BenchmarkInsertRandomRecordsToManticore(b *testing.B) {
	var err error
	var db *sql.DB
	db, err = InitSphinxConnection("9306")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	err = InsertRandomRecordToSphinx(db, recordCount)
	if err != nil {
		log.Println(err)
		return
	}
}

func BenchmarkInsertRandomRecordsToSphinx(b *testing.B) {
	var err error
	var db *sql.DB
	db, err = InitSphinxConnection("9307")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	err = InsertRandomRecordToSphinx(db, recordCount)
	if err != nil {
		log.Println(err)
		return
	}
}

func BenchmarkInsertRandomRecordsToElastic(b *testing.B) {
	var err error
	var client *elastic.Client
	client, err = InitElasticConnection()
	if err != nil {
		log.Println(err)
		return
	}

	err = InsertRandomRecordToElastic(client, recordCount)
	if err != nil {
		log.Println(err)
		return
	}
}

func BenchmarkReadManticoreIds(b *testing.B) {
	var err error
	var db *sql.DB
	db, err = InitSphinxConnection("9306")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	var i int
	for i = 0; i < searchCount; i++ {
		_, err = ReadSphinxIds(db, "MATCH('"+loremIpsumGenerator.Word()+"')")
		if err != nil {
			log.Println(err)
			return
		}
	}

}

func BenchmarkReadSphinxIds(b *testing.B) {
	var err error
	var db *sql.DB
	db, err = InitSphinxConnection("9307")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	var i int
	for i = 0; i < searchCount; i++ {
		_, err = ReadSphinxIds(db, "MATCH('"+loremIpsumGenerator.Word()+"')")
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func BenchmarkReadElasticIds(b *testing.B) {
	var err error
	var client *elastic.Client
	client, err = InitElasticConnection()
	if err != nil {
		log.Println(err)
		return
	}

	var i int
	for i = 0; i < searchCount; i++ {
		_, err = ReadElasticIds(client)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
