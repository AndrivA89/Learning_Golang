/*
	Попробовал работу с reflect'ом

	Добавляю данные со структур User в данные структуры Person
	Создаю слайс из указателей на структуру Person и
	добавляю туда несколько тестовых пользователей

	Вывожу данные до добавления информации со структуры User
	Произвожу добавление данных с помощью рефлексии
	Если поле пустое и есть тег 'default', то берется значение из тега
	Вывожу данные после добавления информации со структуры User

	2-й сценарий:
	Расчет суммы всех транспортных средств
	В функцию суммирования поступает любая структура, которая описывает ТС
	Находим поле со стоимостью и плюсуем
	В случае отсутствия поля или необходимого типа - возвращаем ошибку
*/

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

type Moto struct {
	Model             string
	HorsePower        int
	NumberOfCylinders int
	Cost              float64
}

type Avto struct {
	Model      string
	HorsePower int
	Type       string
	Cost       float64
	Mileage    float64
}

type Vechicle struct {
	Type          string
	HorsePower    int
	Cost          float64
	YearOfRelease int
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

func AddInfo(person *Person, dataToAdd interface{}) {
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
}

func CalculatingTheCost(data []interface{}) (cost float64, err error) {
	var bufErr string
	for i, d := range data {
		if reflect.TypeOf(d).Kind() == reflect.Struct {
			for j := 0; j < reflect.ValueOf(d).NumField(); j++ {
				if reflect.TypeOf(d).Field(j).Name == "Cost" {
					cost += reflect.ValueOf(d).Field(j).Float()
				}
			}
		} else {
			bufErr += "Data " + strconv.Itoa(i) + " NOT a Structure!!!\n"
		}
	}
	err = fmt.Errorf(bufErr)
	return
}

// Сценарий работы со структурами Person и User
func PersonAndUser() {
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

func VechicleAndCost() {
	vechicle := Vechicle{Cost: 5500.11}
	avto := Avto{Cost: 1250.22}
	moto := Moto{Cost: 250.33}
	errorData := []int{12, 13, 14}

	data := make([]interface{}, 0)
	data = append(data, vechicle, avto, moto, errorData, errorData)

	cost, err := CalculatingTheCost(data)
	fmt.Printf("Cost all vechicles equals - %1.2f\n", cost)
	if err != nil {
		fmt.Printf("\nErrors:\n%s", err)
	}
}

func main() {
	PersonAndUser()
	fmt.Printf("\n\n*****************\nExample number 2:\n*****************\n\n")
	VechicleAndCost()
}
