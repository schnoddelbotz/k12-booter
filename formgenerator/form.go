package formgenerator

import (
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Form struct {
	Elements []FormElement
	Complete bool
}
type FormElementType int
type FormElementInputType int
type FormElement struct {
	name          string
	elementType   FormElementType
	inputType     FormElementInputType
	aLink         string
	aURL          string
	selectOptions []FormSelectOption
}

type FormSelectOption struct {
	value    string
	label    string
	selected bool
}

const (
	FT_Unknown FormElementType = iota
	FT_Input
	FT_Select
	FT_HyperLink
	InputType_Text FormElementInputType = iota
	InputType_Number
	InputType_Submit
)

var HTMLElementTypeMap = map[string]FormElementType{
	"input":  FT_Input,
	"select": FT_Select,
	"a":      FT_HyperLink,
}

var HTMLElementInputTypeMap = map[string]FormElementInputType{
	"text":   InputType_Text,
	"number": InputType_Number,
	"submit": InputType_Submit,
}

func ElementTypeHTMLMap() (m map[FormElementType]string) {
	m = make(map[FormElementType]string)
	for k, v := range HTMLElementTypeMap {
		m[v] = k
	}
	return
}

func ElementInputTypeHTMLMap() (m map[FormElementInputType]string) {
	// generics in action, not yet m(, brb ...
	m = make(map[FormElementInputType]string)
	for k, v := range HTMLElementInputTypeMap {
		m[v] = k
	}
	return
}

func (e *FormElement) GetTypeName() string {
	return ElementTypeHTMLMap()[e.elementType]
}

func (e *FormElement) GetInputTypeName() string {
	return ElementInputTypeHTMLMap()[e.inputType]
}

func (f *Form) AddHTMLNodeAsElement(n *html.Node) {
	eType, exists := HTMLElementTypeMap[n.Data]
	if exists {
		elem := &FormElement{}
		elem.elementType = eType
		elem.name = getElementAttributeValue(n.Attr, "name")

		switch eType {
		case FT_Input:
			inputType := getElementAttributeValue(n.Attr, "type")
			iType, exists := HTMLElementInputTypeMap[inputType]
			if !exists {
				log.Panicf("Failed to look up element type %v", inputType)
			}
			elem.inputType = FormElementInputType(iType)
		case FT_Select:
			elem.selectOptions = getFormSelectOptions(n)
		// a.href ...
		default:
			log.Panicf("Can't deal with this element in forms yet: %v", n)
		}

		f.Elements = append(f.Elements, *elem)
	} else {
		log.Printf("Form contains element I cannot work with yet: %s", n.Data)
	}
}

func getElementAttributeValue(attrs []html.Attribute, elementName string) string {
	result := ""
	for _, v := range attrs {
		if v.Key == elementName {
			return v.Val
		}
	}
	// fixme silent failure :/
	return result
}

func elementAttributeExists(attrs []html.Attribute, elementName string) bool {
	// fixme silent failure :/ ... continued / part 2
	for _, v := range attrs {
		if v.Key == elementName {
			return true
		}
	}
	return false
}

func getFormSelectOptions(n *html.Node) []FormSelectOption {
	options := []FormSelectOption{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if c.Data != "option" {
			continue
		}
		option := FormSelectOption{}

		if c.FirstChild != nil {
			option.label = trimNewlinesAndWhitespace(renderNode(c.FirstChild))
		}

		option.value = getElementAttributeValue(c.Attr, "value")

		if elementAttributeExists(c.Attr, "selected") {
			option.selected = true
		}

		options = append(options, option)
	}

	return options
}

func trimNewlinesAndWhitespace(txt string) string {
	txt = strings.TrimSuffix(txt, "\n")
	txt = strings.TrimPrefix(txt, "\n")
	txt = strings.TrimSpace(txt)
	return txt
}
