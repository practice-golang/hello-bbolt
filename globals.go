package main

import (
	"github.com/blevesearch/bleve/v2"
)

const (
	databaseName   = "store.db"
	bucketPeople   = "people"
	bucketPeopleID = "idPeople" // id 증분 관리용 버킷
	indexName      = "index.bleve"
	pageSize       = 4096 // BBolt의 기본 페이지 크기
)

var (
	// listenIP = "0.0.0.0"
	listenIP   = "127.0.0.1"
	listenPORT = "12480"
	listenADDR = listenIP + ":" + listenPORT

	encryptionKey = []byte("mysecretkey12345")
	index         bleve.Index
	db            *EncryptedDB
)
