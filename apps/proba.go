package main

import (
	"log"

	"gopkg.in/yaml.v2"
)

type Item struct {
	A string
	B string
	C string
}

var data = `
- a: a1
  b: b1
  c: c1
- a: a2
  b: b2
  c: c2
- a: a3
  b: b3
  c: c3
`

func main() {

	//list0 := []Item{
	//	{"a1", "b1", "c1"},
	//	{"a2", "b2", "c2"},
	//	{"a3", "b3", "c3"},
	//}
	//
	//data, err := yaml.Marshal(list0)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Print(string(data))

	var list []Item

	err := yaml.Unmarshal([]byte(data), &list)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", list)

}
