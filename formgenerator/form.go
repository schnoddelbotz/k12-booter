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
	id       string // for checkbox and radio ~ labels
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
	InputType_Radio
	InputType_Checkbox
)

var HTMLElementTypeMap = map[string]FormElementType{
	"input":  FT_Input,
	"select": FT_Select,
	"a":      FT_HyperLink,
	"label":  FT_Label,
	"br":     FT_LineBreak,
}

var HTMLElementInputTypeMap = map[string]FormElementInputType{
	"text":     InputType_Text,
	"number":   InputType_Number,
	"submit":   InputType_Submit,
	"radio":    InputType_Radio,
	"checkbox": InputType_Checkbox,
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
	if !exists {
		log.Printf("Form contains element I cannot work with yet: %s", n.Data)
		return
	}
	elem := &FormElement{}
	elem.elementType = eType
	elem.name = getElementAttributeValue(n.Attr, "name")
	elem.id = getElementAttributeValue(n.Attr, "id")
	inputType := getElementAttributeValue(n.Attr, "type")

	switch eType {
	case FT_Input:
		iType, exists := HTMLElementInputTypeMap[inputType]
		if !exists {
			log.Panicf("Unsupported HTML <input> type=%v", inputType)
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

	if eType == FT_Input && inputType == "radio" || inputType == "checkbox" {
		sopt := FormSelectOption{
			value:    getElementAttributeValue(n.Attr, "value"),
			label:    elem.name,
			id:       elem.id,
			selected: elementAttributeExists(n.Attr, "checked"),
		}
		existingId := f.getInternalElementId(elem.name)
		if existingId != nil {
			f.Elements[*existingId].selectOptions = append(f.Elements[*existingId].selectOptions, sopt)
			return
		}
		elem.selectOptions = append(elem.selectOptions, sopt)
	}

	f.Elements = append(f.Elements, *elem)
}

func (f *Form) getInternalElementId(inputName string) *int {
	for id, fdata := range f.Elements {
		if fdata.name == inputName {
			return &id
		}
	}
	return nil
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
