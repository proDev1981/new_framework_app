package ui

import . "app/core"

func App() Element {

	return Box(Args{Name: "root"},
		Title(),
		Controler(),
		ListPerson(35),
	)
}
func algo(){
	fmt.Println("algo")
}