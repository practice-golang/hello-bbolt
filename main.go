package main // import "hello-bbolt"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"go.etcd.io/bbolt"
)

type Person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Birth   string `json:"birth"`   // yyyy-mm-dd
	RegDTTM string `json:"regdttm"` // yyyyMMddhhmmss
}

var (
	db *bbolt.DB
)

func setupDB() error {
	var err error

	db, err = bbolt.Open("store.db", 0666, nil)
	if err != nil {
		panic(err.Error())
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("people"))
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

func addInitialData(db *bbolt.DB) error {
	persons := []Person{
		{ID: 1, Name: "Alice", Birth: "1990-01-01", RegDTTM: "20230725083000"},
		{ID: 2, Name: "Bob", Birth: "1985-05-15", RegDTTM: "20230725093000"},
		{ID: 3, Name: "Carol", Birth: "1992-07-21", RegDTTM: "20230725103000"},
		{ID: 4, Name: "David", Birth: "1980-11-30", RegDTTM: "20230725113000"},
		{ID: 5, Name: "Eve", Birth: "1991-03-12", RegDTTM: "20230725123000"},
		{ID: 6, Name: "Frank", Birth: "1978-04-25", RegDTTM: "20230725133000"},
		{ID: 7, Name: "Grace", Birth: "1995-06-07", RegDTTM: "20230725143000"},
		{ID: 8, Name: "Hank", Birth: "1982-09-16", RegDTTM: "20230725153000"},
		{ID: 9, Name: "Ivy", Birth: "1993-02-18", RegDTTM: "20230725163000"},
		{ID: 10, Name: "Jack", Birth: "1988-10-23", RegDTTM: "20230725173000"},
		{ID: 11, Name: "Kara", Birth: "1994-12-05", RegDTTM: "20230725183000"},
		{ID: 12, Name: "Leo", Birth: "1987-07-09", RegDTTM: "20230725193000"},
		{ID: 13, Name: "Mia", Birth: "1991-05-20", RegDTTM: "20230725203000"},
		{ID: 14, Name: "Nina", Birth: "1983-01-11", RegDTTM: "20230725213000"},
		{ID: 15, Name: "Oscar", Birth: "1992-08-14", RegDTTM: "20230725223000"},
		{ID: 16, Name: "Paul", Birth: "1986-04-23", RegDTTM: "20230725233000"},
		{ID: 17, Name: "Quinn", Birth: "1994-09-30", RegDTTM: "20230726083000"},
		{ID: 18, Name: "Rita", Birth: "1980-12-22", RegDTTM: "20230726093000"},
		{ID: 19, Name: "Steve", Birth: "1990-06-15", RegDTTM: "20230726103000"},
		{ID: 20, Name: "Tina", Birth: "1989-11-28", RegDTTM: "20230726113000"},
		{ID: 21, Name: "Uma", Birth: "1995-01-07", RegDTTM: "20230726123000"},
		{ID: 22, Name: "Vera", Birth: "1987-08-29", RegDTTM: "20230726133000"},
		{ID: 23, Name: "Will", Birth: "1992-04-05", RegDTTM: "20230726143000"},
		{ID: 24, Name: "Xena", Birth: "1985-09-17", RegDTTM: "20230726153000"},
		{ID: 25, Name: "Yara", Birth: "1993-10-11", RegDTTM: "20230726163000"},
		{ID: 26, Name: "Zack", Birth: "1988-06-24", RegDTTM: "20230726173000"},
		{ID: 27, Name: "Amy", Birth: "1991-12-13", RegDTTM: "20230726183000"},
		{ID: 28, Name: "Ben", Birth: "1986-03-01", RegDTTM: "20230726193000"},
		{ID: 29, Name: "Cathy", Birth: "1990-11-21", RegDTTM: "20230726203000"},
		{ID: 30, Name: "Dan", Birth: "1984-05-16", RegDTTM: "20230726213000"},
		{ID: 31, Name: "Ella", Birth: "1992-07-18", RegDTTM: "20230726223000"},
		{ID: 32, Name: "Fred", Birth: "1983-02-25", RegDTTM: "20230726233000"},
		{ID: 33, Name: "Gina", Birth: "1995-03-09", RegDTTM: "20230727083000"},
		{ID: 34, Name: "Holly", Birth: "1989-09-04", RegDTTM: "20230727093000"},
		{ID: 35, Name: "Ian", Birth: "1991-01-27", RegDTTM: "20230727103000"},
		{ID: 36, Name: "Jane", Birth: "1986-07-16", RegDTTM: "20230727113000"},
		{ID: 37, Name: "Ken", Birth: "1994-02-20", RegDTTM: "20230727123000"},
		{ID: 38, Name: "Lena", Birth: "1982-06-09", RegDTTM: "20230727133000"},
		{ID: 39, Name: "Mike", Birth: "1990-11-15", RegDTTM: "20230727143000"},
		{ID: 40, Name: "Nina", Birth: "1985-05-22", RegDTTM: "20230727153000"},
		{ID: 41, Name: "Oliver", Birth: "1991-03-18", RegDTTM: "20230727163000"},
		{ID: 42, Name: "Paula", Birth: "1987-08-31", RegDTTM: "20230727173000"},
		{ID: 43, Name: "Quincy", Birth: "1992-04-16", RegDTTM: "20230727183000"},
		{ID: 44, Name: "Riley", Birth: "1983-12-13", RegDTTM: "20230727193000"},
		{ID: 45, Name: "Sam", Birth: "1995-09-23", RegDTTM: "20230727203000"},
		{ID: 46, Name: "Tina", Birth: "1986-07-20", RegDTTM: "20230727213000"},
		{ID: 47, Name: "Ursula", Birth: "1991-10-30", RegDTTM: "20230727223000"},
		{ID: 48, Name: "Vince", Birth: "1984-03-25", RegDTTM: "20230727233000"},
		{ID: 49, Name: "Wendy", Birth: "1995-12-09", RegDTTM: "20230728083000"},
		{ID: 50, Name: "Xander", Birth: "1988-05-04", RegDTTM: "20230728093000"},
		{ID: 51, Name: "Yvonne", Birth: "1992-01-15", RegDTTM: "20230728103000"},
		{ID: 52, Name: "Zane", Birth: "1983-08-27", RegDTTM: "20230728113000"},
		{ID: 53, Name: "Alice", Birth: "1990-12-30", RegDTTM: "20230728123000"},
		{ID: 54, Name: "Bob", Birth: "1987-11-02", RegDTTM: "20230728133000"},
		{ID: 55, Name: "Carol", Birth: "1995-03-20", RegDTTM: "20230728143000"},
		{ID: 56, Name: "David", Birth: "1982-01-13", RegDTTM: "20230728153000"},
		{ID: 57, Name: "Eve", Birth: "1994-08-07", RegDTTM: "20230728163000"},
		{ID: 58, Name: "Frank", Birth: "1986-06-19", RegDTTM: "20230728173000"},
		{ID: 59, Name: "Grace", Birth: "1991-02-10", RegDTTM: "20230728183000"},
		{ID: 60, Name: "Hank", Birth: "1985-12-15", RegDTTM: "20230728193000"},
		{ID: 61, Name: "Ivy", Birth: "1993-06-04", RegDTTM: "20230729103000"},
		{ID: 62, Name: "Jack", Birth: "1989-04-25", RegDTTM: "20230729113000"},
		{ID: 63, Name: "Kara", Birth: "1992-11-15", RegDTTM: "20230729123000"},
		{ID: 64, Name: "Leo", Birth: "1984-03-05", RegDTTM: "20230729133000"},
		{ID: 65, Name: "Mia", Birth: "1995-05-25", RegDTTM: "20230729143000"},
		{ID: 66, Name: "Nina", Birth: "1987-09-09", RegDTTM: "20230729153000"},
		{ID: 67, Name: "Oscar", Birth: "1991-07-22", RegDTTM: "20230729163000"},
		{ID: 68, Name: "Paul", Birth: "1986-11-13", RegDTTM: "20230729173000"},
		{ID: 69, Name: "Quinn", Birth: "1994-01-31", RegDTTM: "20230729183000"},
		{ID: 70, Name: "Rita", Birth: "1985-10-19", RegDTTM: "20230729193000"},
		{ID: 71, Name: "Steve", Birth: "1993-12-10", RegDTTM: "20230730103000"},
		{ID: 72, Name: "Tina", Birth: "1988-02-20", RegDTTM: "20230730113000"},
		{ID: 73, Name: "Uma", Birth: "1994-06-11", RegDTTM: "20230730123000"},
		{ID: 74, Name: "Vera", Birth: "1989-01-15", RegDTTM: "20230730133000"},
		{ID: 75, Name: "Will", Birth: "1992-08-03", RegDTTM: "20230730143000"},
		{ID: 76, Name: "Xena", Birth: "1984-05-17", RegDTTM: "20230730153000"},
		{ID: 77, Name: "Yara", Birth: "1990-03-21", RegDTTM: "20230730163000"},
		{ID: 78, Name: "Zack", Birth: "1986-12-30", RegDTTM: "20230730173000"},
		{ID: 79, Name: "Amy", Birth: "1995-04-14", RegDTTM: "20230730183000"},
		{ID: 80, Name: "Ben", Birth: "1987-07-05", RegDTTM: "20230730193000"},
	}

	err := db.Update(func(tx *bbolt.Tx) error {
		personBucket := tx.Bucket([]byte("people"))
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

		for _, p := range persons {

			data, err := json.Marshal(p)
			if err != nil {
				return err
			}
			key := fmt.Sprintf("%d", p.ID)
			if err := personBucket.Put([]byte(key), data); err != nil {
				return err
			}

			// name 인덱싱
			if err := nameIndexBucket.Put([]byte(p.Name), []byte(key)); err != nil {
				return err
			}

			// birth 인덱싱
			birthKey := p.Birth + fmt.Sprintf("%d", p.ID)
			if err := birthIndexBucket.Put([]byte(birthKey), []byte(key)); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func addPerson(db *bbolt.DB, p Person) error {
	return db.Update(func(tx *bbolt.Tx) error {
		personBucket := tx.Bucket([]byte("people"))
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

		if personBucket.Get([]byte(fmt.Sprintf("%d", p.ID))) != nil {
			return fmt.Errorf("person with ID %d already exists", p.ID)
		}

		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		if err := personBucket.Put([]byte(fmt.Sprintf("%d", p.ID)), data); err != nil {
			return err
		}

		// name 인덱싱
		if err := nameIndexBucket.Put([]byte(p.Name), []byte(fmt.Sprintf("%d", p.ID))); err != nil {
			return err
		}
		// birth 인덱싱
		birthKey := p.Birth + fmt.Sprintf("%d", p.ID)
		if err := birthIndexBucket.Put([]byte(birthKey), []byte(fmt.Sprintf("%d", p.ID))); err != nil {
			return err
		}

		return nil
	})
}

func updatePerson(db *bbolt.DB, id int, newPerson Person) error {
	return db.Update(func(tx *bbolt.Tx) error {
		personBucket := tx.Bucket([]byte("people"))
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

		oldPersonData := personBucket.Get([]byte(fmt.Sprintf("%d", id)))
		if oldPersonData == nil {
			return fmt.Errorf("person with ID %d not found", id)
		}

		data, err := json.Marshal(newPerson)
		if err != nil {
			return err
		}
		if err := personBucket.Put([]byte(fmt.Sprintf("%d", id)), data); err != nil {
			return err
		}

		// name 인덱싱
		oldPerson := Person{}
		if err := json.Unmarshal(oldPersonData, &oldPerson); err != nil {
			return err
		}
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
}

func deletePerson(db *bbolt.DB, id int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		personBucket := tx.Bucket([]byte("people"))
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

		// name 인덱싱
		p := Person{}
		if err := json.Unmarshal(personData, &p); err != nil {
			return err
		}
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
}

func searchByName(name string, db *bbolt.DB) {
	err := db.View(func(tx *bbolt.Tx) error {
		nameIndexBucket := tx.Bucket([]byte("name_index"))
		if nameIndexBucket == nil {
			return fmt.Errorf("name index bucket not found")
		}

		personID := nameIndexBucket.Get([]byte(name))
		if personID == nil {
			fmt.Printf("Person with name %s not found\n", name)
			return nil
		}

		personBucket := tx.Bucket([]byte("people"))
		if personBucket == nil {
			return fmt.Errorf("person bucket not found")
		}

		personData := personBucket.Get(personID)
		if personData == nil {
			fmt.Printf("Person data with ID %s not found\n", personID)
			return nil
		}

		var p Person
		if err := json.Unmarshal(personData, &p); err != nil {
			return err
		}

		fmt.Printf("Found person: %+v\n", p)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func searchByBirthRange(startDateStr, endDateStr string, db *bbolt.DB) {
	err := db.View(func(tx *bbolt.Tx) error {
		birthIndexBucket := tx.Bucket([]byte("birth_index"))
		if birthIndexBucket == nil {
			return fmt.Errorf("birth index bucket not found")
		}

		personBucket := tx.Bucket([]byte("people"))
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

				// fmt.Printf("Person in birth range: ID=%s\n", id)
				fmt.Printf("Person in birth range: %+v\n", p)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error

	err = setupDB()
	if err != nil {
		panic(err.Error())
	}

	err = addInitialData(db)
	if err != nil {
		panic(err.Error())
	}

	addPerson(db, Person{ID: 81, Name: "Zara", Birth: "1989-03-22", RegDTTM: "20230731183000"})       // 1 추가
	updatePerson(db, 1, Person{ID: 1, Name: "Alice", Birth: "1991-01-01", RegDTTM: "20230731083000"}) // 수정

	deletePerson(db, 2)                                // 삭제
	searchByName("Alice", db)                          // 이름 검색
	searchByBirthRange("1990-01-01", "1992-12-31", db) // 날짜 범위 검색

	db.Close()
}
