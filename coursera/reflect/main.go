package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

var countGuest = 0

type User struct {
	ID    int
	Login string `default:"guest"`
	Posts int    `default:"empty"`
}

type Person struct {
	User User
	Name string
	Age  int
}

func GetUser(id int) *User {
	user := new(User)
	if id > 10 {
		user.Login = ""
		user.Posts = 0
	} else {
		user.Login = "User_" + strconv.Itoa(id)
		rand.Seed(time.Now().UnixNano())
		user.Posts = rand.Intn(101)
	}
	user.ID = id
	return user
}

func AddInfo(person *Person, dataToAdd interface{}) (err error) {
	data := reflect.ValueOf(dataToAdd).Elem()
	typeField := reflect.TypeOf(dataToAdd).Elem()
	for i := 0; i < data.NumField(); i++ {
		switch typeField.Field(i).Name {
		case "ID":
			person.User.ID = int(data.Field(i).Int())
		case "Login":
			if login := data.Field(i).String(); login != "" {
				person.User.Login = login
			} else {
				if defaultValue := data.Type().Field(i).Tag.Get("default"); defaultValue != "" {
					countGuest++
					person.User.Login = defaultValue + "_" + strconv.Itoa(countGuest)
				}
			}
		case "Posts":
			person.User.Posts = int(data.Field(i).Int())
		}
	}

	return
}

func main() {
	persons := []*Person{}

	test1 := &Person{Name: "Andrew", Age: 30}
	test2 := &Person{Name: "Alex", Age: 22}
	test3 := &Person{Name: "Fedor", Age: 54}
	test4 := &Person{Name: "Sergey", Age: 13}
	test5 := &Person{Name: "Marina", Age: 41}

	persons = append(persons, test1, test2, test3, test4, test5)

	fmt.Println("\n***************\nPersons before:\n***************")
	for i := 0; i < len(persons); i++ {
		fmt.Printf("\nName - %s, Age - %d\nUserID - %d, Login - %s, Posts - %d\n",
			persons[i].Name, persons[i].Age,
			persons[i].User.ID, persons[i].User.Login, persons[i].User.Posts)
	}

	AddInfo(persons[0], GetUser(1))
	AddInfo(persons[1], GetUser(5))
	AddInfo(persons[2], GetUser(11))
	AddInfo(persons[3], GetUser(3))
	AddInfo(persons[4], GetUser(99))

	fmt.Println("\n**************\nPersons after:\n**************")
	for i := 0; i < len(persons); i++ {
		fmt.Printf("\nName - %s, Age - %d\nUserID - %d, Login - %s, Posts - %d\n",
			persons[i].Name, persons[i].Age,
			persons[i].User.ID, persons[i].User.Login, persons[i].User.Posts)
	}
}
