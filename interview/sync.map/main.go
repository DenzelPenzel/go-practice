package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map
	m.Store("address", map[string]string{"name1": "Vasya", "name2": "Petya"})
	v, _ := m.Load("address")
	// invalid operation: v["name1"] (
	// type interface {} does not support indexing) Because func (m *Map) Store(key interface{}, value interface{})
	// so v type is interface {}, a type assertion is needed here

	fmt.Println(v.(map[string]string)["name1"])
}
