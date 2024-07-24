package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	bolt "go.etcd.io/bbolt"
)

// getNextID는 자동 증분 ID를 가져오거나 생성합니다.
func getNextID(db *EncryptedDB) (int, error) {
	var nextID int

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketPeopleID))
		if err != nil {
			return err
		}

		// 마지막 ID를 취득
		lastID := b.Get([]byte("lastID"))
		if lastID == nil {
			nextID = 1 // 시작 ID를 1로 설정
		} else {
			// 마지막 ID 정수 변환하고 +1
			lastIDInt, err := strconv.Atoi(string(lastID))
			if err != nil {
				return err
			}
			nextID = lastIDInt + 1
		}

		return b.Put([]byte("lastID"), []byte(strconv.Itoa(nextID))) // 다음 ID를 버킷에 저장
	})

	if err != nil {
		return -1, err
	}

	return nextID, nil
}

func AddPerson(p Person) error {
	// 자동 증분 ID 생성
	id, err := getNextID(db)
	if err != nil {
		return err
	}
	p.ID = id
	p.RegDTTM = time.Now().Format("20060102150405")

	// 데이터베이스에 추가
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

func GetAllPersons() ([]Person, error) {
	var persons []Person

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketPeople))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var person Person
			decryptedValue, err := decryptPage(v, encryptionKey)
			if err != nil {
				return err
			}
			err = json.Unmarshal(decryptedValue, &person)
			if err != nil {
				return err
			}
			persons = append(persons, person)
			return nil
		})
	})

	return persons, err
}

func GetPerson(id string) (Person, error) {
	var person Person
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketPeople))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		encryptedValue := b.Get([]byte(id))
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

func SearchPeople(query string) {
	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Printf("Error searching: %v", err)
		return
	}

	fmt.Printf("Found %d matches for query '%s':\n", searchResult.Total, query)
	for _, hit := range searchResult.Hits {
		person, err := GetPerson(hit.ID)
		if err != nil {
			log.Printf("Error getting person: %v", err)
			continue
		}
		fmt.Printf("- %s (Age: %d, Gender: %s)\n", person.Name, person.Age, person.Gender)
	}
}

func SearchPeopleByDate(startDate, endDate time.Time) {
	query := bleve.NewDateRangeQuery(startDate, endDate)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Printf("Error searching: %v", err)
		return
	}

	fmt.Printf("Found %d matches for date range %v to %v:\n", searchResult.Total, startDate, endDate)
	for _, hit := range searchResult.Hits {
		person, err := GetPerson(hit.ID)
		if err != nil {
			log.Printf("Error getting person: %v", err)
			continue
		}

		regdttm, _ := time.Parse("20060102150405", person.RegDTTM)

		fmt.Printf("- ID: %d, Name: %s, Age: %d, Gender: %s, RegDTTM: %s\n",
			person.ID, person.Name, person.Age, person.Gender, regdttm.Format("2006-01-02 15:04:05"))
	}
}
func SearchPeopleByAge(startAge, endAge int) {
	var startAgeF float64 = float64(startAge)
	var endAgeF float64 = float64(endAge)

	query := bleve.NewNumericRangeQuery(&startAgeF, &endAgeF)
	query.SetField("Age")
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Printf("Error searching: %v", err)
		return
	}

	fmt.Printf("Found %d matches for age range %d to %d:\n", searchResult.Total, startAge, endAge)
	for _, hit := range searchResult.Hits {
		person, err := GetPerson(hit.ID)
		if err != nil {
			log.Printf("Error getting person: %v", err)
			continue
		}

		regdttm, _ := time.Parse("20060102150405", person.RegDTTM)

		fmt.Printf("- ID: %d, Name: %s, Age: %d, Gender: %s, RegDTTM: %s\n",
			person.ID, person.Name, person.Age, person.Gender, regdttm.Format("2006-01-02 15:04:05"))
	}
}
