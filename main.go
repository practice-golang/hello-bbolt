package main // import "hello-bbolt"

import (
	"fmt"

	"go.etcd.io/bbolt"
)

var xchacha *Xchacha

func setupInitialData(db *bbolt.DB) {
	persons := []Person{
		{Name: "Alice", Gender: "Female", Birth: "1990-01-01", RegDTTM: "20230725083000"},
		{Name: "Bob", Gender: "Male", Birth: "1985-05-15", RegDTTM: "20230725093000"},
		{Name: "Carol", Gender: "Female", Birth: "1992-07-21", RegDTTM: "20230725103000"},
		{Name: "David", Gender: "Male", Birth: "1980-11-30", RegDTTM: "20230725113000"},
		{Name: "Eve", Gender: "Female", Birth: "1991-03-12", RegDTTM: "20230725123000"},
		{Name: "Frank", Gender: "Male", Birth: "1978-04-25", RegDTTM: "20230725133000"},
		{Name: "Grace", Gender: "Female", Birth: "1995-06-07", RegDTTM: "20230725143000"},
		{Name: "Hank", Gender: "Male", Birth: "1982-09-16", RegDTTM: "20230725153000"},
		{Name: "Ivy", Gender: "Female", Birth: "1993-02-18", RegDTTM: "20230725163000"},
		{Name: "Jack", Gender: "Male", Birth: "1988-10-23", RegDTTM: "20230725173000"},
		{Name: "Kara", Gender: "Female", Birth: "1994-12-05", RegDTTM: "20230725183000"},
		{Name: "Leo", Gender: "Male", Birth: "1987-07-09", RegDTTM: "20230725193000"},
		{Name: "Mia", Gender: "Female", Birth: "1991-05-20", RegDTTM: "20230725203000"},
		{Name: "Nina", Gender: "Female", Birth: "1983-01-11", RegDTTM: "20230725213000"},
		{Name: "Oscar", Gender: "Male", Birth: "1992-08-14", RegDTTM: "20230725223000"},
		{Name: "Paul", Gender: "Male", Birth: "1986-04-23", RegDTTM: "20230725233000"},
		{Name: "Quinn", Gender: "Male", Birth: "1994-09-30", RegDTTM: "20230726083000"},
		{Name: "Rita", Gender: "Female", Birth: "1980-12-22", RegDTTM: "20230726093000"},
		{Name: "Steve", Gender: "Male", Birth: "1990-06-15", RegDTTM: "20230726103000"},
		{Name: "Tina", Gender: "Female", Birth: "1989-11-28", RegDTTM: "20230726113000"},
		{Name: "Uma", Gender: "Female", Birth: "1995-01-07", RegDTTM: "20230726123000"},
		{Name: "Vera", Gender: "Female", Birth: "1987-08-29", RegDTTM: "20230726133000"},
		{Name: "Will", Gender: "Male", Birth: "1992-04-05", RegDTTM: "20230726143000"},
		{Name: "Xena", Gender: "Female", Birth: "1985-09-17", RegDTTM: "20230726153000"},
		{Name: "Yara", Gender: "Female", Birth: "1993-10-11", RegDTTM: "20230726163000"},
		{Name: "Zack", Gender: "Male", Birth: "1988-06-24", RegDTTM: "20230726173000"},
		{Name: "Amy", Gender: "Female", Birth: "1991-12-13", RegDTTM: "20230726183000"},
		{Name: "Ben", Gender: "Male", Birth: "1986-03-01", RegDTTM: "20230726193000"},
		{Name: "Cathy", Gender: "Female", Birth: "1990-11-21", RegDTTM: "20230726203000"},
		{Name: "Dan", Gender: "Male", Birth: "1984-05-16", RegDTTM: "20230726213000"},
		{Name: "Ella", Gender: "Female", Birth: "1992-07-18", RegDTTM: "20230726223000"},
		{Name: "Fred", Gender: "Male", Birth: "1983-02-25", RegDTTM: "20230726233000"},
		{Name: "Gina", Gender: "Female", Birth: "1995-03-09", RegDTTM: "20230727083000"},
		{Name: "Holly", Gender: "Male", Birth: "1989-09-04", RegDTTM: "20230727093000"},
		{Name: "Ian", Gender: "Male", Birth: "1991-01-27", RegDTTM: "20230727103000"},
		{Name: "Jane", Gender: "Female", Birth: "1986-07-16", RegDTTM: "20230727113000"},
		{Name: "Ken", Gender: "Male", Birth: "1994-02-20", RegDTTM: "20230727123000"},
		{Name: "Lena", Gender: "Female", Birth: "1982-06-09", RegDTTM: "20230727133000"},
		{Name: "Mike", Gender: "Male", Birth: "1990-11-15", RegDTTM: "20230727143000"},
		{Name: "Oliver", Gender: "Male", Birth: "1991-03-18", RegDTTM: "20230727163000"},
		{Name: "Paula", Gender: "Male", Birth: "1987-08-31", RegDTTM: "20230727173000"},
		{Name: "Quincy", Gender: "Female", Birth: "1992-04-16", RegDTTM: "20230727183000"},
		{Name: "Riley", Gender: "Male", Birth: "1983-12-13", RegDTTM: "20230727193000"},
		{Name: "Sam", Gender: "Male", Birth: "1995-09-23", RegDTTM: "20230727203000"},
		{Name: "Ursula", Gender: "Female", Birth: "1991-10-30", RegDTTM: "20230727223000"},
		{Name: "Vince", Gender: "Male", Birth: "1984-03-25", RegDTTM: "20230727233000"},
		{Name: "Wendy", Gender: "Female", Birth: "1995-12-09", RegDTTM: "20230728083000"},
		{Name: "Xander", Gender: "Male", Birth: "1988-05-04", RegDTTM: "20230728093000"},
		{Name: "Yvonne", Gender: "Female", Birth: "1992-01-15", RegDTTM: "20230728103000"},
		{Name: "Zane", Gender: "Female", Birth: "1983-08-27", RegDTTM: "20230728113000"},
	}

	for _, p := range persons {
		err := addPerson(db, p)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func main() {
	var err error

	err = setupDB()
	if err != nil {
		panic(err.Error())
	}

	xchacha, err = SetupXchacha()
	if err != nil {
		panic(err.Error())
	}

	setupInitialData(db)

	addPerson(db, Person{Name: "Zara", Birth: "1989-03-22", RegDTTM: "20230731183000"})        // 1 추가
	updatePerson(db, 1, Person{Name: "Alice", Birth: "1991-01-02", RegDTTM: "20230731083000"}) // 수정

	deletePerson(db, 2)                    // 삭제
	found, err := searchByName("Mike", db) // 이름 검색
	if err != nil {
		panic(err.Error())
	}
	if found.RegDTTM != "" {
		fmt.Printf("Person search by name: %+v\n", found)
	}

	persons, err := SearchByBirthRange("1988-01-01", "1990-01-01", db) // 날짜 범위 검색
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Person in birth range:")
	for _, p := range persons {
		fmt.Printf("%+v\n", p)
	}

	db.Close()
}
