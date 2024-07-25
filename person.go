package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	bolt "go.etcd.io/bbolt"
)

func AddPerson(p Person) error {
	id, err := getNextID(db) // ID 생성
	if err != nil {
		return err
	}

	p.ID = id
	p.RegDTTM = time.Now().Format("20060102150405")

	// db에 추가
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketPeople))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		encryptedValue, err := encryptPage(buf, encryptionKey)
		if err != nil {
			return err
		}

		err = b.Put(itob(p.ID), encryptedValue)
		if err != nil {
			return err
		}

		return index.Index(fmt.Sprintf("%d", p.ID), p)
	})

	return err
}

func DeletePerson(id int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketPeople))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// 삭제 대상
		idBytes := itob(id)
		existingValue := b.Get(idBytes)
		if existingValue == nil {
			return fmt.Errorf("person with ID %d not found", id)
		}

		// db에서 삭제
		err := b.Delete(idBytes)
		if err != nil {
			return err
		}

		// bleve 인덱스에서 삭제
		return index.Delete(fmt.Sprintf("%d", id))
	})

	return err
}

func UpdatePerson(id int, updatedPerson Person) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketPeople))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// 기존 정보
		idBytes := itob(id)
		existingValue := b.Get(idBytes)
		if existingValue == nil {
			return fmt.Errorf("person with ID %d not found", id)
		}

		// 기존 데이터
		var existingPerson Person
		decryptedValue, err := decryptPage(existingValue, encryptionKey)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(decryptedValue, &existingPerson); err != nil {
			return err
		}

		// ID, RegDTTM 유지
		updatedPerson.ID = id
		updatedPerson.RegDTTM = existingPerson.RegDTTM

		// 수정된 정보를 JSON 변환
		updatedValue, err := json.Marshal(updatedPerson)
		if err != nil {
			return err
		}

		// 값 암호화
		encryptedValue, err := encryptPage(updatedValue, encryptionKey)
		if err != nil {
			return err
		}

		// db 저장
		err = b.Put(idBytes, encryptedValue)
		if err != nil {
			return err
		}

		// bleve 인덱스 업데이트
		return index.Index(fmt.Sprintf("%d", id), updatedPerson)
	})
}

func GetPerson(id int) (Person, error) {
	var person Person

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketPeople))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		idBytes := itob(id)
		encryptedValue := b.Get(idBytes)
		if encryptedValue == nil {
			return fmt.Errorf("person not found")
		}

		decryptedValue, err := decryptPage(encryptedValue, encryptionKey)
		if err != nil {
			return err
		}

		return json.Unmarshal(decryptedValue, &person)
	})

	return person, err
}

func GetPersonList(search PersonSearch, limit int) ([]Person, error) {
	var err error

	que := bleve.NewBooleanQuery()

	// 이름: 일부 일치
	if search.Name != "" {
		nameQuery := bleve.NewWildcardQuery("*" + search.Name + "*")
		nameQuery.SetField("name")
		que.AddMust(nameQuery)
	}

	// 성별: 전부 일치
	if search.Gender != "" {
		genderQuery := bleve.NewMatchQuery(search.Gender)
		genderQuery.SetField("gender")
		que.AddMust(genderQuery)
	}

	// from, to: 범위
	if search.From != "" && search.To != "" {
		var fromTime, toTime time.Time

		if search.From != "" {
			fromTime, err = time.Parse("2006-01-02", search.From)
			if err != nil {
				return nil, fmt.Errorf("invalid From date format: %v", err)
			}
		}

		if search.To != "" {
			toTime, err = time.Parse("2006-01-02", search.To)
			if err != nil {
				return nil, fmt.Errorf("invalid To date format: %v", err)
			}
			toTime = toTime.Add(24*time.Hour - time.Second) // To 시간 - 23:59:59
		}

		dateQuery := bleve.NewDateRangeQuery(fromTime, toTime)
		dateQuery.SetField("birth")
		que.AddMust(dateQuery)
	}

	// Boolean 쿼리는 최소 1개 이상 조건 필요
	if que.Validate() != nil {
		wildcardQuery := bleve.NewWildcardQuery("*")
		que.AddMust(wildcardQuery)
	}

	// 검색 요청 생성
	searchRequest := bleve.NewSearchRequest(que)
	searchRequest.Size = limit

	if search.Sort == "DESC" {
		searchRequest.SortBy([]string{"-birth", "_score"}) // SORT DESC
	} else {
		searchRequest.SortBy([]string{"birth", "time", "_score"}) // SORT ASC
	}

	// 검색 실행
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search error: %v", err)
	}

	// 결과 처리
	var persons []Person
	for _, hit := range searchResults.Hits {
		id, _ := strconv.ParseInt(hit.ID, 10, 64)
		person, err := GetPerson(int(id))
		if err != nil {
			return nil, err
		}

		persons = append(persons, person)
	}

	return persons, nil
}
