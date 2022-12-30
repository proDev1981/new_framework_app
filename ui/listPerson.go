package ui

import . "app/core"
import . "app/model"

func ListPerson(max int) Element {

	list := NewState("list", DataMO)

	return List(Args{Class: "map", State: list, Max: max},
		Label(Args{
			Class: "fecha",
			Type:  "span",
			Value: "{{.Fecha}}",
		}),
		Label(Args{
			Class: "name",
			Type:  "span",
			Value: "{{.Name}}",
		}),
		Label(Args{
			Class: "obra",
			Type:  "span",
			Value: "{{.Obra}}",
		}),
	)
}
