package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.etcd.io/bbolt"
)

type Person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Birth   string `json:"birth"`   // yyyy-mm-dd
	RegDTTM string `json:"regdttm"` // yyyyMMddhhmmss
}

var db *bbolt.DB

// // name_index list
// db.View(func(tx *bbolt.Tx) error {
// 	bucket := tx.Bucket([]byte("name_index"))
// 	if bucket == nil {
// 		return fmt.Errorf("bucket %s not found", "name_index")
// 	}

// 	err := bucket.ForEach(func(k, v []byte) error {
// 		name, err := decrypt(k)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return err
// 		}

// 		fmt.Printf("key=%s, value=%s\n", name, v)
// 		return nil
// 	})

// 	return err
// })

func setupDB() error {
	var err error

	db, err = bbolt.Open("store.db", 0666, nil)
	if err != nil {
		panic(err.Error())
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("xchacha"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("person"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("person_id"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("name_index"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("birth_index"))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func KeyExists(bucketName string, key []byte) bool {
	var exists bool

	db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return nil // 버킷이 없으면 키도 없는 것으로 간주
		}

		value := bucket.Get(key)
		exists = value != nil

		return nil
	})

	return exists
}

// 자동 증분 ID 생성
func getNextID(db *bbolt.DB) (int, error) {
	var nextID int

	err := db.Update(func(tx *bbolt.Tx) error {
		personBucket := tx.Bucket([]byte("person_id"))
		if personBucket == nil {
			return fmt.Errorf("person bucket not found")
		}

		lastID := personBucket.Get([]byte("lastID")) // 마지막 ID 취득
		if lastID == nil {
			nextID = 1 // 시작 ID 1로 설정
		} else {
			lastIDInt, err := strconv.Atoi(string(lastID)) // 마지막 ID 정수 변환하고 +1
			if err != nil {
				return err
			}
			nextID = lastIDInt + 1
		}

		return personBucket.Put([]byte("lastID"), []byte(strconv.Itoa(nextID))) // 다음 ID를 버킷에 저장
	})

	if err != nil {
		return -1, err
	}

	return nextID, nil
}

func addPerson(db *bbolt.DB, p Person) error {
	var err error

	id, err := getNextID(db) // ID 생성
	if err != nil {
		return err
	}

	p.ID = id

	err = db.Update(func(tx *bbolt.Tx) error {
		personBucket := tx.Bucket([]byte("person"))
		if personBucket == nil {
			return fmt.Errorf("person bucket not found")
		}

		if personBucket.Get([]byte(fmt.Sprintf("%d", p.ID))) != nil {
			return fmt.Errorf("person with ID %d already exists", p.ID)
		}

		decName := p.Name
		encName, err := encrypt(strings.TrimSpace(decName))
		if err != nil {
			return fmt.Errorf("failed to encrypt")
		}
		p.Name = encName

		data, err := json.Marshal(p)
		if err != nil {
			return err
		}

		nameIndexBucket := tx.Bucket([]byte("name_index"))
		if nameIndexBucket == nil {
			return fmt.Errorf("name index bucket not found")
		}
		birthIndexBucket := tx.Bucket([]byte("birth_index"))
		if birthIndexBucket == nil {
			return fmt.Errorf("birth index bucket not found")
		}

		// name 중복 체크
		if nameIndexBucket.Get([]byte(encName)) != nil {
			return fmt.Errorf("name %s already exists", decName)
		}

		if err := personBucket.Put([]byte(fmt.Sprintf("%d", p.ID)), data); err != nil {
			return err
		}

		// name 인덱싱
		if err := nameIndexBucket.Put([]byte(encName), []byte(fmt.Sprintf("%d", p.ID))); err != nil {
			return err
		}
		// birth 인덱싱
		birthKey := p.Birth + fmt.Sprintf("%d", p.ID)
		if err := birthIndexBucket.Put([]byte(birthKey), []byte(fmt.Sprintf("%d", p.ID))); err != nil {
			return err
		}

		return nil
	})

	return err
}

func updatePerson(db *bbolt.DB, id int, newPerson Person) error {
	err := db.Update(func(tx *bbolt.Tx) error {
		personBucket := tx.Bucket([]byte("person"))
		if personBucket == nil {
			return fmt.Errorf("person bucket not found")
		}

		oldPerson := Person{}
		oldPersonData := personBucket.Get([]byte(fmt.Sprintf("%d", id)))
		if oldPersonData == nil {
			return fmt.Errorf("person with ID %d not found", id)
		}
		if err := json.Unmarshal(oldPersonData, &oldPerson); err != nil {
			return err
		}

		encName, err := encrypt(newPerson.Name)
		if err != nil {
			return err
		}
		newPerson.Name = encName
		newPerson.ID = oldPerson.ID
		newPerson.RegDTTM = oldPerson.RegDTTM

		data, err := json.Marshal(newPerson)
		if err != nil {
			return err
		}
		if err := personBucket.Put([]byte(fmt.Sprintf("%d", id)), data); err != nil {
			return err
		}

		nameIndexBucket := tx.Bucket([]byte("name_index"))
		if nameIndexBucket == nil {
			return fmt.Errorf("name index bucket not found")
		}
		birthIndexBucket := tx.Bucket([]byte("birth_index"))
		if birthIndexBucket == nil {
			return fmt.Errorf("birth index bucket not found")
		}

		// name 인덱싱
		if oldPerson.Name != newPerson.Name {
			if err := nameIndexBucket.Delete([]byte(oldPerson.Name)); err != nil {
				return err
			}
			if err := nameIndexBucket.Put([]byte(newPerson.Name), []byte(fmt.Sprintf("%d", id))); err != nil {
				return err
			}
		}

		// birth 인덱싱
		if oldPerson.Birth != newPerson.Birth {
			oldBirthKey := oldPerson.Birth + fmt.Sprintf("%d", id)
			if err := birthIndexBucket.Delete([]byte(oldBirthKey)); err != nil {
				return err
			}
			newBirthKey := newPerson.Birth + fmt.Sprintf("%d", id)
			if err := birthIndexBucket.Put([]byte(newBirthKey), []byte(fmt.Sprintf("%d", id))); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func deletePerson(db *bbolt.DB, id int) error {
	err := db.Update(func(tx *bbolt.Tx) error {
		personBucket := tx.Bucket([]byte("person"))
		if personBucket == nil {
			return fmt.Errorf("person bucket not found")
		}

		nameIndexBucket := tx.Bucket([]byte("name_index"))
		if nameIndexBucket == nil {
			return fmt.Errorf("name index bucket not found")
		}
		birthIndexBucket := tx.Bucket([]byte("birth_index"))
		if birthIndexBucket == nil {
			return fmt.Errorf("birth index bucket not found")
		}

		personData := personBucket.Get([]byte(fmt.Sprintf("%d", id)))
		if personData == nil {
			return fmt.Errorf("person with ID %d not found", id)
		}

		if err := personBucket.Delete([]byte(fmt.Sprintf("%d", id))); err != nil {
			return err
		}

		p := Person{}
		if err := json.Unmarshal(personData, &p); err != nil {
			return err
		}
		name, err := encrypt(strings.TrimSpace(p.Name))
		if err != nil {
			return fmt.Errorf("failed to encrypt")
		}
		p.Name = name

		// name 인덱싱
		if err := nameIndexBucket.Delete([]byte(p.Name)); err != nil {
			return err
		}
		// birth 인덱싱
		birthKey := p.Birth + fmt.Sprintf("%d", id)
		if err := birthIndexBucket.Delete([]byte(birthKey)); err != nil {
			return err
		}

		return nil
	})

	return err
}

func searchByName(name string, db *bbolt.DB) (Person, error) {
	var err error
	var result Person

	encName, err := encrypt(name)
	if err != nil {
		return Person{}, fmt.Errorf("failed to encrypt")
	}

	err = db.View(func(tx *bbolt.Tx) error {
		nameIndexBucket := tx.Bucket([]byte("name_index"))
		if nameIndexBucket == nil {
			return fmt.Errorf("name index bucket not found")
		}

		personID := nameIndexBucket.Get([]byte(encName))
		if personID == nil {
			fmt.Printf("Person with name %s not found\n", name)
			return nil
		}

		personBucket := tx.Bucket([]byte("person"))
		if personBucket == nil {
			return fmt.Errorf("person bucket not found")
		}

		personData := personBucket.Get(personID)
		if personData == nil {
			fmt.Printf("Person data with ID %s not found\n", personID)
			return nil
		}

		if err := json.Unmarshal(personData, &result); err != nil {
			return err
		}

		result.Name, err = decrypt([]byte(result.Name))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return Person{}, err
	}

	return result, nil
}

func SearchByBirthRange(startDateStr, endDateStr string, db *bbolt.DB) ([]Person, error) {
	var err error
	var result []Person

	err = db.View(func(tx *bbolt.Tx) error {
		birthIndexBucket := tx.Bucket([]byte("birth_index"))
		if birthIndexBucket == nil {
			return fmt.Errorf("birth index bucket not found")
		}

		personBucket := tx.Bucket([]byte("person"))
		if personBucket == nil {
			return fmt.Errorf("person bucket not found")
		}

		cursor := birthIndexBucket.Cursor()
		for k, id := cursor.Seek([]byte(startDateStr)); k != nil && bytes.Compare(k, []byte(endDateStr)) <= 0; k, id = cursor.Next() {
			birthDate := string(k)
			if birthDate >= startDateStr && birthDate <= endDateStr {
				personData := personBucket.Get(id)
				if personData == nil {
					continue
				}
				var p Person
				if err := json.Unmarshal(personData, &p); err != nil {
					return err
				}

				if p.Name == "" {
					continue
				}

				decName, err := decrypt([]byte(p.Name))
				if err != nil {
					return err
				}
				p.Name = decName

				result = append(result, p)
			}
		}

		return nil
	})

	return result, err
}
