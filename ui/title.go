package ui

import . "app/core"

func Title() Element {

	tittle := NewState("tittle", "Mano de obra")

	return Box(Args{
		Class: "container-title"},
		Label(Args{
			Class: "title",
			State: tittle,
			Value: "{{.state}}"}),
	)
}
