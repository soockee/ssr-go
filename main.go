package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/soockee/ssr-go/components"
)

func main() {

	component := components.Hello("hello")

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
