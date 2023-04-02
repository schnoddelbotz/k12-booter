package formgenerator

import (
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Form struct {
	Elements []FormElement
	labels   map[string]string
	Complete bool
}
type FormElementType int
type FormElementInputType int
type FormElement struct {
	name          string
	id            string
	labelText     string // any <label> text found
	elementType   FormElementType
	inputType     FormElementInputType
	aLink         string
	aURL          string
	selectOptions []FormSelectOption
	text          string
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
	FT_LineBreak
	FT_Label
	FT_TextElement
	InputType_Text FormElementInputType = iota
	InputType_Number
	InputType_Submit
)

var HTMLElementTypeMap = map[string]FormElementType{
	"input":  FT_Input,
	"select": FT_Select,
	"a":      FT_HyperLink,
	"label":  FT_Label,
	"br":     FT_LineBreak,
}

var HTMLElementInputTypeMap = map[string]FormElementInputType{
	"text":   InputType_Text,
	"number": InputType_Number,
	"submit": InputType_Submit,
}

func NewForm() *Form {
	f := &Form{}
	f.labels = make(map[string]string)
	return f
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
	// todo: move to output_html.go
	eType, exists := HTMLElementTypeMap[n.Data]
	if exists {
		elem := &FormElement{}
		elem.elementType = eType
		elem.name = getElementAttributeValue(n.Attr, "name")
		elem.id = getElementAttributeValue(n.Attr, "id")

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
		case FT_LineBreak:
		case FT_Label:
			forid := getElementAttributeValue(n.Attr, "for")
			labelText := renderNode(n.FirstChild)
			f.labels[forid] = labelText
			// only append to f.Labels, not f.Elements - thus early:
			return
		// a.href ...
		default:
			log.Panicf("Can't deal with this element in forms yet: %v", n)
		}

		f.Elements = append(f.Elements, *elem)
	} else {
		log.Printf("Form contains element I cannot work with yet: %s", n.Data)
	}
}

func (f *Form) SetLabelTextPerID() {
	// merge labels stored separately into corresponding elements
	for eid, fe := range f.Elements {
		if fe.id != "" {
			label, ok := f.labels[fe.id]
			if ok {
				f.Elements[eid].labelText = label
			}
		}
	}
}

func (f *Form) AddTextAsElement(text string) {
	f.Elements = append(f.Elements, FormElement{elementType: FT_TextElement, text: text})
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
