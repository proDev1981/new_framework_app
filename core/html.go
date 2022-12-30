package core

import "fmt"
import "strings"
import "log"
import "reflect"
import "os"
import "assets"
import "encoding/json"
import "github.com/gorilla/websocket"

var dom Element

// init
func init() {
	log.SetFlags(log.Lshortfile)
}

type Html struct {
	server   *Server
	dom      Element
	conn     *websocket.Conn
	provider *Provider
}

var html *Html

// contructor
func NewHtml() *Html {
	if html == nil {
		return &Html{server: NewServer(), provider: newProvider()}
	}
	return html
}

// builder page html
func HtmlBuild(p *page) Motor {
	p.SetMotorRender(NewHtml())
	if _, err := os.Stat("./src"); err != nil {
		log.Println("No exist folder => src ")
		log.Println("Creating this..")
		if err = os.Mkdir("./src", os.ModePerm); err != nil {
			panic(err)
		}
	}
	file, err := os.Create("./src/index.html")
	if err != nil {
		panic(err)
	}
	file.Write([]byte(p.Render()))
	return p.motorRender
}

// add listener in struct html
// @ params
// parent is id of element html
// listener is map[type evenet]function event
func (h *Html) AddEventListener(parent string, listener Listener) {
	for types, call := range listener {
		if h.server.GetSocket() == nil {
			h.server.SetInitialEvents(parent, types, call)
		} else {
			name := types + parent
			h.server.GetSocket().SendAll(ComposeBind(name))
			h.server.GetSocket().SendAll(ComposeEventListener(parent, types, name))
		}
	}
}

// getter server
func (h *Html) NewServer() *Server {
	return h.server
}

// getter server
func (h *Html) GetServer() *Server {
	return h.server
}

// render map element
func (h *Html) RenderMap(e Element) string {
	var body string
	var header string
	var footer string

	header += fmt.Sprint("<", e.Tag(), argsToHTml(e.Args()), ">")
	footer += fmt.Sprint("</", e.Tag(), ">")

	if e.Args().State != nil {
		types := fmt.Sprintf("%T", e.Args().State)
		if types == "*core.State" {
			data := e.Args().State.Get()
			Dvalue := reflect.ValueOf(data)
			lenArrayData := Dvalue.Len()
			if lenArrayData > e.Args().Max {
				lenArrayData = e.Args().Max
			}
			for i := 0; i < lenArrayData; i++ {
				Dstruct := Dvalue.Index(i)
				argsContainer := e.Args()
				argsContainer.Class = e.Args().Class + "-item"
				body += fmt.Sprint("<", e.Tag(), argsToHTml(argsContainer), ">")
				for _, item := range e.Children() {
					item.setParent(e)
					lenFieldsInStruct := Dstruct.NumField()
					innerHtml := item.Args().Value
					for i := 0; i < lenFieldsInStruct; i++ {
						value := fmt.Sprint(Dstruct.Field(i))
						key := fmt.Sprint(Dstruct.Type().Field(i).Name)
						innerHtml = strings.ReplaceAll(innerHtml, "{{."+key+"}}", value)
					}
					body += fmt.Sprint("<", item.Tag(), argsToHTml(item.Args()), ">", innerHtml, "</", item.Tag(), ">")
				}
				body += fmt.Sprint("</", e.Tag(), ">")
			}

		}
	}
	return fmt.Sprint(header, body, footer)
}

// render element for html
func (h *Html) RenderElement(e Element) (res string) {
	if e.GetSubType() == "list" {
		state := e.Args().State
		e.getProvider().AddState(state)
		state.Add(e)
		res = h.RenderMap(e)
	} else {
		var value string
		if e.Args().State != nil {
			state := e.Args().State
			e.getProvider().AddState(state)
			state.Add(e)
			value = replaceState(e.Args(), "")
		} else {
			value = e.Args().Value
		}
		res = fmt.Sprint("<", e.Tag(), argsToHTml(e.Args()), ">", value)
		for _, item := range e.Children() {
			item.setParent(e)
			item.SetMotorRender(h)
			res += item.render()
			if !eventsIsEmpty(item.Args().Events) {
				h.AddEventListener(item.Args().id, item.Args().Events)
			}
			if item.Args().Link != nil {
				event := Listener{"change": func(e *Event) {
					// reparar link no se graba en bien
					//*item.Args().Link = <-e.Target().GetAttribute("value").Await()
				}}
				h.AddEventListener(item.Args().id, event)
			}
		}
		res += fmt.Sprint("</", e.Tag(), ">")
	}
	return
}

// render element for html
func (h *Html) RenderPage(p *page) (res string) {
	var body string
	var headers string

	res = fmt.Sprint("<!DOCTYPEhtml><html>")
	for _, item := range p.Children() {
		item.SetMotorRender(p.motorRender)
		if item.Tag() != "style" && item.Tag() != "head" {
			if item.Args().Name == "root" {
				dom = item
				h.dom = item
			}
			body += item.render()
		} else {
			headers += item.render()
		}
	}
	res += fmt.Sprint(headers, "<body>", body, "</body></html>")
	return res
}

// search element by query
func (h *Html) RootSelector(query string) Element {
	return selector(dom, query)
}

// search element by query
func (h *Html) Selector(e Element, query string) Element {
	return selector(e, query)
}

// search element by query
func (h *Html) SelectorAll(e Element, query string) []Element {
	return selectorAll(e, query)
}

// search element by query
func selector(e Element, query string) Element {
	if e != nil {
		return search(query, e)
	}
	return nil

}

// search element by query
func selectorAll(e Element, query string) (res []Element) {
	if e != nil {
		res = append(res, searchAll(query, e)...)
	}
	return

}

// setter conn
func (h *Html) SetConn(conn *websocket.Conn) {
	h.conn = conn
}

// getter conn
func (h *Html) Conn() *websocket.Conn {
	return h.conn
}

// update element in dom of client
func (h *Html) Update(e Element) {
	compose := ComposeEval("document.getElementById('%s').outerHTML =`%s`",
		e.Args().id, e.render(),
	)
	h.conn.WriteMessage(1, []byte(compose))
}

// create new objetc in js
func (h *Html) NewObject(e Element, name string, value any) *PROMISE {
	p := NewPromise(h)
	val := strings.ReplaceAll(string(assets.Try(json.Marshal(value))), `"`, "'")
	createObject := ComposeEval(`window.%s = new Object`, name)
	fillObject := ComposeEvalAnsResponse(p.id, `%s=%s`, name, val)
	h.Conn().WriteMessage(1, []byte(createObject))
	h.Conn().WriteMessage(1, []byte(fillObject))
	return p
}

// set attribute of element in dom html
func (h *Html) SetAttribute(e Element, name, value string) *PROMISE {
	var cmd string
	if name == "value" {
		if e.GetSubType() == "button" {
			name = "innerHTML"
		}
		cmd = `document.getElementById('%s').%s='%s'`
	} else {
		cmd = `document.getElementById('%s').setAttribute('%s','%s')`
	}
	p := NewPromise(h)
	setAttribute := ComposeEval(
		cmd,
		e.Args().id,
		name,
		value,
	)
	h.Conn().WriteMessage(1, []byte(setAttribute))
	return p
}

// geter attribute of element in dom html
func (h *Html) GetAttribute(e Element, name string) *PROMISE {
	p := NewPromise(h)
	var getAttribute string
	if name == "value" {
		if e.GetSubType() == "button" {
			getAttribute = ComposeEvalAnsResponse(p.id,
				`document.getElementById('%s').innerHTML`,
				e.Args().id,
			)
		} else {
			getAttribute = ComposeEvalAnsResponse(p.id,
				`document.getElementById('%s').value`,
				e.Args().id,
			)
		}
	} else {
		getAttribute = ComposeEvalAnsResponse(p.id,
			`document.getElementById('%s').getAttribute('%s')`,
			e.Args().id,
			name,
		)
	}
	h.Conn().WriteMessage(1, []byte(getAttribute))
	return p
}

// log string in js console
func (h *Html) Log(e Element, value ...string) {
	log := ComposeEval(`console.log('%s')`, value)
	h.Conn().WriteMessage(1, []byte(log))
}

// alert in js client
func (h *Html) Alert(e Element, value ...string) {
	alert := ComposeEval(`alert('%s')`, value)
	h.Conn().WriteMessage(1, []byte(alert))
}

// get target
func Target(e *Event) Element {
	ele := selector(dom, "#"+e.Id)
	ele.MotorRender().SetConn(e.Client)
	return ele
}

// return map string string values input in element
func (h *Html) GetData(e Element) map[string]string {

	children := e.SelectorAll("input")
	patter := make(map[string]string)
	if len(children) > 0 {
		for _, item := range children {
			patter[item.Args().Class] = <-item.GetAttribute("value").Await()
		}
	}
	return patter
}

// return provider
func (h *Html) getProvider() *Provider {
	return h.provider
}

// return state from map provider
func (h *Html) GetState(name string) *State {
	return h.provider.GetState(name)
}

// delete value element
func (h *Html) Reset(e Element) Element {
	e.SetAttribute("value", "")
	return e
}

// focus element
func (h *Html) Focus(e Element) Element {
	focus := ComposeEval(`document.getElementById('%s').focus()`, e.Args().id)
	h.Conn().WriteMessage(1, []byte(focus))
	return e
}

// get innerhtml element
func (h *Html) GetInner(e Element) *PROMISE {
	p := NewPromise(h)
	getInnerHtml := ComposeEvalAnsResponse(p.id,
		"document.getElementById('%s').innerHTML",
		e.Args().id,
	)
	h.Conn().WriteMessage(1, []byte(getInnerHtml))
	return p
}

// set innerHTML element
func (h *Html) SetInner(e Element, value string) {
	setInnerHTml := ComposeEval(
		"document.getElementById('%s').innerHTML = '%s'",
		e.Args().id,
		value,
	)
	h.Conn().WriteMessage(1, []byte(setInnerHTml))
}

// get value element
func (h *Html) GetValue(e Element) *PROMISE {
	return h.GetAttribute(e, "value")
}

// set value element
func (h *Html) SetValue(e Element, value string) {
	h.SetAttribute(e, "value", value)
}
