package main

import (
	"encoding/binary"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

func ResponseJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// 자동 증분 ID 생성
func getNextID(db *EncryptedDB) (int, error) {
	var nextID int

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketPeopleID))
		if err != nil {
			return err
		}

		// 마지막 ID 취득
		lastID := b.Get([]byte("lastID"))
		if lastID == nil {
			nextID = 1 // 시작 ID 1로 설정
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

// ID를 바이트 슬라이스로 변환
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))

	return b
}

func InitData() {
	// 데이터 추가
	var people = []Person{
		{Name: "Alice", Age: 30, Gender: "Female", Birth: "1994-05-01"},
		{Name: "Bob", Age: 25, Gender: "Male", Birth: "1998-03-15"},
		{Name: "Charlie", Age: 35, Gender: "Male", Birth: "1988-07-22"},
		{Name: "David", Age: 40, Gender: "Male", Birth: "1983-11-30"},
		{Name: "Eve", Age: 22, Gender: "Female", Birth: "2001-06-10"},
		{Name: "Frank", Age: 28, Gender: "Male", Birth: "1995-08-05"},
		{Name: "Grace", Age: 29, Gender: "Female", Birth: "1994-01-20"},
		{Name: "Hank", Age: 33, Gender: "Male", Birth: "1990-12-12"},
		{Name: "Ivy", Age: 31, Gender: "Female", Birth: "1993-04-25"},
		{Name: "Jack", Age: 27, Gender: "Male", Birth: "1996-09-18"},
		{Name: "Karen", Age: 24, Gender: "Female", Birth: "1999-10-14"},
		{Name: "Leo", Age: 32, Gender: "Male", Birth: "1991-02-22"},
		{Name: "Mona", Age: 26, Gender: "Female", Birth: "1997-11-08"},
		{Name: "Nick", Age: 37, Gender: "Male", Birth: "1986-07-30"},
		{Name: "Olivia", Age: 23, Gender: "Female", Birth: "2000-05-17"},
		{Name: "Paul", Age: 39, Gender: "Male", Birth: "1984-03-25"},
		{Name: "Quinn", Age: 34, Gender: "Female", Birth: "1989-06-12"},
		{Name: "Randy", Age: 41, Gender: "Male", Birth: "1982-01-09"},
		{Name: "Sara", Age: 27, Gender: "Female", Birth: "1996-03-23"},
		{Name: "Tom", Age: 30, Gender: "Male", Birth: "1993-08-01"},
		{Name: "Ursula", Age: 35, Gender: "Female", Birth: "1988-12-18"},
		{Name: "Victor", Age: 29, Gender: "Male", Birth: "1994-10-05"},
		{Name: "Wendy", Age: 24, Gender: "Female", Birth: "1999-02-28"},
		{Name: "Xander", Age: 38, Gender: "Male", Birth: "1985-04-17"},
		{Name: "Yara", Age: 28, Gender: "Female", Birth: "1995-09-09"},
		{Name: "Zach", Age: 31, Gender: "Male", Birth: "1992-11-21"},
		{Name: "Ava", Age: 26, Gender: "Female", Birth: "1997-07-02"},
		{Name: "Brian", Age: 40, Gender: "Male", Birth: "1983-12-07"},
		{Name: "Clara", Age: 22, Gender: "Female", Birth: "2001-04-14"},
		{Name: "Derek", Age: 33, Gender: "Male", Birth: "1990-09-16"},
		{Name: "Ella", Age: 27, Gender: "Female", Birth: "1996-05-28"},
		{Name: "Finn", Age: 32, Gender: "Male", Birth: "1991-06-23"},
		{Name: "Gina", Age: 30, Gender: "Female", Birth: "1993-03-10"},
		{Name: "Henry", Age: 35, Gender: "Male", Birth: "1988-08-14"},
		{Name: "Iris", Age: 29, Gender: "Female", Birth: "1994-02-06"},
		{Name: "James", Age: 28, Gender: "Male", Birth: "1995-12-19"},
		{Name: "Kara", Age: 36, Gender: "Female", Birth: "1987-11-02"},
		{Name: "Leo", Age: 31, Gender: "Male", Birth: "1993-04-17"},
		{Name: "Mia", Age: 24, Gender: "Female", Birth: "1999-08-23"},
		{Name: "Nate", Age: 37, Gender: "Male", Birth: "1986-07-12"},
		{Name: "Olive", Age: 33, Gender: "Female", Birth: "1991-12-09"},
		{Name: "Paul", Age: 26, Gender: "Male", Birth: "1997-11-29"},
		{Name: "Quincy", Age: 29, Gender: "Female", Birth: "1994-09-15"},
		{Name: "Rachel", Age: 32, Gender: "Female", Birth: "1992-05-30"},
		{Name: "Sam", Age: 28, Gender: "Male", Birth: "1995-01-03"},
		{Name: "Tina", Age: 34, Gender: "Female", Birth: "1989-07-16"},
		{Name: "Ulysses", Age: 31, Gender: "Male", Birth: "1992-06-07"},
		{Name: "Vera", Age: 27, Gender: "Female", Birth: "1996-04-22"},
		{Name: "Will", Age: 30, Gender: "Male", Birth: "1993-09-28"},
		{Name: "Xena", Age: 26, Gender: "Female", Birth: "1997-10-10"},
	}

	regDate := time.Now()

	for i, p := range people {
		regdttm := regDate.AddDate(0, 0, i)
		p.RegDTTM = regdttm.Format("20060102150405")
		AddPerson(p)
	}
}
