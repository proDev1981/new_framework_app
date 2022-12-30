package model

import (
	"csv"
	"strings"

	. "extend"
)

type ManoObra struct {
	Name  string
	Fecha string
	Obra  string
}
type ListManoObra []ManoObra

// / covert [][]string to ManoObra struct
func CsvToManoObra(data [][]string) ListManoObra {
	var this []ManoObra
	for _, item := range data[1:] {
		if len(item) > 2 {
			item[2] = item[2][:len(item[2])-1] // borro en caracter controlM de unix
			this = append(this, ManoObra{Name: item[1], Fecha: item[2], Obra: item[0]})
		}
	}
	return this
}

// add new element in person
func (p ListManoObra) Add(name, date, job string) ListManoObra {
	return append(p, ManoObra{name, date, job})
}

// filter elements by query in persons
func (p ListManoObra) Filter(value string) (res ListManoObra) {
	value = strings.ToLower(value)
	for _, item := range p {
		if strings.Contains(strings.ToLower(item.Name), value) ||
			strings.Contains(strings.ToLower(String(item.Fecha)), value) ||
			strings.Contains(strings.ToLower(item.Obra), value) {
			res = append(res, item)
		}
	}
	return
}

// filter AND elements by querys in persons
func (p ListManoObra) FilterOR(values ...string) (res ListManoObra) {
	for _, value := range values {
		res = append(res, p.Filter(value)...)
	}
	return
}

// filter OR element by querys in persons
func (p ListManoObra) FilterAND(values ...string) (res ListManoObra) {
	res = p
	for _, value := range values {
		res = res.Filter(value)
	}
	return
}

// return reverse
func (p ListManoObra) Reverse() (res ListManoObra) {
	max := len(p) - 1
	for index := range p {
		res = append(res, p[max-index])
	}
	return
}

var DataMO = CsvToManoObra(csv.Open(`G:/Mi unidad/DB/src/operarios.db`).Get()).Reverse()
