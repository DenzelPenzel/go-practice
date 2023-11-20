package main

import "fmt"

type Student struct {
	Name string
	Age  int
}

func main() {
	test1()
	test2()
	test3()
	test4()
	test5()
}

func test1() {
	fmt.Println("========")

	kv := map[string]Student{"vasya": {Age: 19}}

	// kv["vasya"].Age = 22 // error, invalid operation, map in golang return two values, first - value, second - whether the map exists key
	fmt.Println("test1 ->", kv)

	var s = []Student{{Age: 19}}
	s[0].Age = 22 // valid operation
	fmt.Println("test1 ->", s, s[0].Age)

	// var str string = nil // error, strings in golang cannot be assigned nil and compared with nil
	// if str == nil {
	//	str = "default"
	//}
}

// ==========================================================================================================
// Pointer Type
func test2() {
	fmt.Println("\n========")

	mp := make(map[string]*Student)

	list := []Student{
		{Name: "vasya", Age: 10},
		{Name: "petya", Age: 12},
		{Name: "denis", Age: 11},
	}

	// each loop will copy the value in the collection
	// therefore, the last value will be stored in the last map
	for _, item := range list {
		mp[item.Name] = &item
	}

	fmt.Println("test2 ->", mp) // map[denis:0xc000010078 petya:0xc000010078 vasya:0xc000010078]
}

// ==========================================================================================================
// no concept of inheritance in the golang language, only combination, no virtual methods, and no overloading
type People struct{}

func (p *People) ShowA() {
	fmt.Println("showA")
	p.ShowB()
}

func (p *People) ShowB() {
	fmt.Println("showB")
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher showB")
}

func test3() {
	fmt.Println("\n========")

	t := Teacher{}
	t.ShowA() // print: showA, showB
}

// ==========================================================================================================
type Base interface {
	Login(string) string
}

type App struct{}

// method has a pointer receiver (*App)
func (app *App) Login(password string) (res string) {
	if password == "pass" {
		res = "login"
	} else {
		res = "sorry wrong password"
	}
	return
}

func test4() {
	fmt.Println("\n========")
	var user Base = &App{} // change to the ref pointer, App{} should't work

	fmt.Println(user.Login("pass"))
}

// ==========================================================================================================
type People2 interface {
	Show()
}

type Student2 struct{}

func (s *Student2) Show() {
}

func live() People2 {
	var stu *Student2
	return stu
}

func test5() {
	fmt.Println("\n========")
	if live() == nil {
		fmt.Println("hello")
	} else {
		fmt.Println("world") // return it
	}
}
