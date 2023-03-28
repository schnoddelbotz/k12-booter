package formgenerator

import (
	"fmt"
	"log"
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

func ReadFormParamsFromFile(filename string) *Form {
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
	f.complete = true
	return f
}

func CreateForm(f *Form) {
	log.Printf("TODO: use template to produce cui boilerplate code: %v", f)
	// should wire hotkeys to void functions to be filled with live ...
}

func CreateFormAsNeeded(byQuery bool, byFilename string) bool {
	form := &Form{}
	if byQuery {
		form = QueryFormParamsFromUser(byQuery)
	}
	if byFilename != "" {
		form = ReadFormParamsFromFile(byFilename)
	}
	if form.complete {
		CreateForm(form)
		return true
	}
	return false
}
