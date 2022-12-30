package model

import "strings"
import . "extend"

type Person struct {
	Name string
	Age  int
}
type Operarios []Person

// add new element in person
func (p Operarios) Add(name string, age int) Operarios {
	return append(p, Person{name, age})
}

// filter elements by query in persons
func (p Operarios) Filter(value string) (res Operarios) {
	for _, item := range p {
		if strings.Contains(item.Name, value) || strings.Contains(String(item.Age), value) {
			res = append(res, item)
		}
	}
	return
}

// filter AND elements by querys in persons
func (p Operarios) FilterOR(values ...string) (res Operarios) {
	for _, value := range values {
		res = append(res, p.Filter(value)...)
	}
	return
}

// filter OR element by querys in persons
func (p Operarios) FilterAND(values ...string) (res Operarios) {
	res = p
	for _, value := range values {
		res = res.Filter(value)
	}
	return
}

// database simulator
var DataOperarios Operarios = Operarios{
	{
		Name: "alberto",
		Age:  42,
	},
	{
		Name: "paco",
		Age:  39,
	},
	{
		Name: "sergio",
		Age:  43,
	},
}
