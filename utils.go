package main

import (
	"encoding/binary"
	"encoding/json"
	"net/http"
	"time"
)

func ResponseJSON(w http.ResponseWriter, statusCode int, message map[string]string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

func InitData() {
	// 데이터 추가
	var people = []Person{
		{Name: "Alice", Age: 30, Gender: "Female"},
		{Name: "Bob", Age: 25, Gender: "Male"},
		{Name: "Charlie", Age: 35, Gender: "Male"},
		{Name: "David", Age: 40, Gender: "Male"},
		{Name: "Eve", Age: 22, Gender: "Female"},
		{Name: "Frank", Age: 28, Gender: "Male"},
		{Name: "Grace", Age: 29, Gender: "Female"},
		{Name: "Hank", Age: 33, Gender: "Male"},
		{Name: "Ivy", Age: 31, Gender: "Female"},
		{Name: "Jack", Age: 27, Gender: "Male"},
		{Name: "Karen", Age: 24, Gender: "Female"},
		{Name: "Leo", Age: 32, Gender: "Male"},
		{Name: "Mona", Age: 26, Gender: "Female"},
		{Name: "Nick", Age: 37, Gender: "Male"},
		{Name: "Olivia", Age: 23, Gender: "Female"},
		{Name: "Paul", Age: 39, Gender: "Male"},
		{Name: "Quinn", Age: 34, Gender: "Female"},
		{Name: "Randy", Age: 41, Gender: "Male"},
		{Name: "Sara", Age: 27, Gender: "Female"},
		{Name: "Tom", Age: 30, Gender: "Male"},
		{Name: "Ursula", Age: 35, Gender: "Female"},
		{Name: "Victor", Age: 29, Gender: "Male"},
		{Name: "Wendy", Age: 24, Gender: "Female"},
		{Name: "Xander", Age: 38, Gender: "Male"},
		{Name: "Yara", Age: 28, Gender: "Female"},
		{Name: "Zach", Age: 31, Gender: "Male"},
		{Name: "Ava", Age: 26, Gender: "Female"},
		{Name: "Brian", Age: 40, Gender: "Male"},
		{Name: "Clara", Age: 22, Gender: "Female"},
		{Name: "Derek", Age: 33, Gender: "Male"},
		{Name: "Ella", Age: 27, Gender: "Female"},
		{Name: "Finn", Age: 32, Gender: "Male"},
		{Name: "Gina", Age: 30, Gender: "Female"},
		{Name: "Henry", Age: 35, Gender: "Male"},
		{Name: "Iris", Age: 29, Gender: "Female"},
		{Name: "James", Age: 28, Gender: "Male"},
		{Name: "Kara", Age: 36, Gender: "Female"},
		{Name: "Leo", Age: 31, Gender: "Male"},
		{Name: "Mia", Age: 24, Gender: "Female"},
		{Name: "Nate", Age: 37, Gender: "Male"},
		{Name: "Olive", Age: 33, Gender: "Female"},
		{Name: "Paul", Age: 26, Gender: "Male"},
		{Name: "Quincy", Age: 29, Gender: "Female"},
		{Name: "Rachel", Age: 32, Gender: "Female"},
		{Name: "Sam", Age: 28, Gender: "Male"},
		{Name: "Tina", Age: 34, Gender: "Female"},
		{Name: "Ulysses", Age: 31, Gender: "Male"},
		{Name: "Vera", Age: 27, Gender: "Female"},
		{Name: "Will", Age: 30, Gender: "Male"},
		{Name: "Xena", Age: 26, Gender: "Female"},
	}

	regDate := time.Now()

	for i, p := range people {
		regdttm := regDate.AddDate(0, 0, i)
		p.RegDTTM = regdttm.Format("20060102150405")
		AddPerson(p)
	}
}

// ID를 바이트 슬라이스로 변환
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
