package core

var p *page

type page struct {
	children    []Element
	motorRender Motor
}

// constructor struct page
func Page(children ...Element) *page {
	if p == nil {
		p = &page{children: children}
	}
	return p
}

// render page
func (p *page) Render() (res string) {
	return p.motorRender.RenderPage(p)
}

// getter children page
func (p *page) Children() []Element {
	return p.children
}

// setter motor render
func (p *page) SetMotorRender(m Motor) {
	p.motorRender = m
}

// getter events listeners
func (p *page) GetMotorRender() Motor {
	return p.motorRender
}

// falta por implementar
func Styles(path string) *Ele {
	return &Ele{tag: "style"}
}
func Link(args Args) *Ele {
	args.Rel = "stylesheet"
	return &Ele{tag: "link", args: args}
}
func Meta(args Args) *Ele {
	return &Ele{tag: "meta", args: args}
}
func PageTitle(args Args) *Ele {
	return &Ele{tag: "title", args: args}
}
func Header(children ...Element) *Ele {
	return &Ele{tag: "head", children: children}
}
func Script(args Args) *Ele {
	return &Ele{tag: "script", args: args}
}
