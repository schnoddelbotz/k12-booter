package formgenerator

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/html"
)

type Form struct {
	elements []FormElement
	complete bool
}

type FormElement struct {
	labelText   string
	elementType FormElementType
}

type FormElementType int

const (
	TypeMenuItem int = iota
	TypePureText
	TypeNewLine
	TypeClosingElement
)

var formNode *html.Node

func QueryFormParamsFromUser(doit bool) *Form {
	// dialog-driven, add elems, ask for filename to write to or stdout?
	if !doit {
		return nil
	}
	form := &Form{complete: false}
	fElement := &FormElement{}
	var input string
	for !form.complete {
		input = ""
		fmt.Print("Enter text: ")
		fmt.Scanln(&input)
		if input == "" {
			form.complete = true
		} else {
			fElement.labelText = input
			fElement.elementType = FormElementType(TypeMenuItem)
			form.elements = append(form.elements, *fElement)
		}
	}
	return form
}

func ReadFormFromHTMLFile(filename string) *Form {
	if filename == "" {
		return nil
	}
	f := &Form{}
	// read dumb simple ascii file, example, feels like markdown?
	// # Main screen
	// 1. Basic local PC k12-booter settings
	// 2. Basic school settings
	// 3. IT Inventory and orders - Hardware
	// 4. IT Inventory and orders - Software
	// 5. IT Configuration management & deployment
	// 6. IT Monitoring, alerting and telemetry
	// 7. Settings report & system documentation
	// 8. Introduction and access to services

	data, err := os.Open(filename)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}
	defer data.Close()
	doc, err := html.Parse(data)
	if err != nil {
		log.Panicf("don't panic. html parser error. return err...")
	}

	zf(doc)
	if formNode == nil {
		log.Print(renderNode(doc))
		log.Panicf("bad html file (no form node found): %v", err)
	}

	htmlform, err := html.ParseFragment(data, formNode)
	if err != nil {
		log.Fatalf("parseFragment failed on form: %v", err)
	}
	for x, r := range htmlform {
		log.Printf("parseFragment found: %v -> %v", x, r.Data)
	}

	n := formNode
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		log.Printf(" C= %+v ", c)
		log.Print(renderNode(c))
	}
	f.complete = true
	return f
}

func zf(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "form" {
		log.Printf("Found form. Form Attributes: %+v", n.Attr)
		log.Print(renderNode(n))
		formNode = n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		zf(c)
	}
}

func CreateFormGO(f *Form) {
	log.Printf("TODO: use template to produce cui boilerplate code: %v", f)
	// should wire hotkeys to void functions to be filled with live ...
}

func CreateFormHTML(f *Form) {
	log.Printf("TODO: use template to produce HTML boilerplate code: %v", f)
	// can be used to diff against imported html form ...
}

func CreateFormAsNeeded(byQuery bool, byFilename string) bool {
	form := &Form{}
	if byQuery {
		form = QueryFormParamsFromUser(byQuery)
	}
	if byFilename != "" {
		form = ReadFormFromHTMLFile(byFilename)
	}
	if form.complete {
		// HTML for now to debug import into *Form, Go later
		// This should reflect the original input form structure. Test...
		CreateFormHTML(form)
		return true
	}
	return false
}

func renderNode(n *html.Node) string {
	// from: https://zetcode.com/golang/net-html/
	var buf bytes.Buffer
	w := io.Writer(&buf)

	err := html.Render(w, n)

	if err != nil {
		return ""
	}

	return buf.String()
}
