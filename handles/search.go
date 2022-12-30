package handles

import . "app/core"
import . "app/model"
import "strings"

// handle btn search
func Search(e *Event) {
	this := e.Target()
	list := this.GetState("list")
	switch <-this.GetValue().Await() {
	case "go":
		search(this, list)

	case "\u21E6":
		reset(this, list)

	}
}

// action search
func search(e Element, list *State) {
	e.SetValue("\u21E6")
	boxSearchValue := <-e.RootSelector(".input-search").GetValue().Await()
	switch {
	case strings.Contains(boxSearchValue, "&&"):
		{
			slice := strings.Split(boxSearchValue, "&&")
			list.Set(DataMO.FilterAND(slice...))
		}
	case strings.Contains(boxSearchValue, "||"):
		{
			slice := strings.Split(boxSearchValue, "||")
			list.Set(DataMO.FilterOR(slice...))
		}
	default:
		list.Set(DataMO.Filter(boxSearchValue))
	}
}

// action reset
func reset(e Element, list *State) {
	e.SetValue("go")
	list.Set(DataMO)
	e.RootSelector(".input-search").Reset().Focus()

}
