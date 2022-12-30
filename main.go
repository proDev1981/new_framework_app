package main

import . "app/core"
import "app/ui"

func main() {

	pwa := HtmlBuild(
		Page(
			Header(
				Link(Args{Href: "./styles.css"}),
				Meta(Args{Charset: "UTF-8"}),
				PageTitle(Args{Value: "Data Base"}),
			),
			Script(Args{Src: "./index.js"}),
			ui.App(),
		))
	pwa.NewServer().AndSocket().Listen()
}
