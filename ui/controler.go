package ui

import . "app/core"
import "app/handles"

// component
func Controler() Element {

	return Box(Args{Class: "box-search"},
		Input(Args{
			Class: "input-search"},
		),
		Button(Args{
			Class:  "btn filter",
			Value:  "go",
			Events: Listener{"click": handles.Search}},
		),
		Button(Args{
			Class:  "btn tidy",
			Value:  "≚",
			Events: Listener{"click": handles.Tidy}},
		),
		Button(Args{
			Class:  "btn add",
			Value:  "+",
			Events: Listener{"click": handles.Add}},
		),
		Button(Args{
			Class:  "btn next",
			Value:  "\u21E8",
			Events: Listener{"click": handles.Next}},
		),
	)
}
