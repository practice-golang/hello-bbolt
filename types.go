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
	RegDTTM string // 등록 날짜: yyyyMMddhhmmss
}
