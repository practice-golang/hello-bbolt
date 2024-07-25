package main

import (
	bolt "go.etcd.io/bbolt"
)

type EncryptedDB struct {
	*bolt.DB
}

type Person struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Birth   string `json:"birth"` // yyyy-mm-dd
	RegDTTM string // 등록 날짜: yyyyMMddhhmmss
}

type PersonSearch struct {
	Name   string `json:"name,omitempty"`   // 이름 검색 (부분 일치)
	Gender string `json:"gender,omitempty"` // 성별 검색 (정확 일치)
	From   string `json:"from,omitempty"`   // 생일 범위 시작
	To     string `json:"to,omitempty"`     // 생일 범위 끝
	Sort   string `json:"sort,omitempty"`   // 정렬
}
