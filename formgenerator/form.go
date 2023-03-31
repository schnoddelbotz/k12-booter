package formgenerator

import (
	"log"

	"golang.org/x/net/html"
)

type Form struct {
	Elements []FormElement
	Complete bool
}

type FormElement struct {
	labelText   string
	elementType FormElementType
}

type FormElementType int

const (
	FT_Input FormElementType = iota
	FT_Select
	FT_HyperLink
)

var HTMLElementTypeMap = map[string]FormElementType{
	"input":  FT_Input,
	"select": FT_Select,
}

func ElementTypeHTMLMap() (m map[FormElementType]string) {
	for k, v := range HTMLElementTypeMap {
		m[v] = k
	}
	return
}

func (f *Form) AddHTMLNodeAsElement(n *html.Node) {
	log.Printf("Append node %s to Form size now/before: %d", n.Data, len(f.Elements))
	eType, exists := HTMLElementTypeMap[n.Data]
	if exists {
		log.Printf("Add n type %s (internal %d)", n.Data, eType)
		f.Elements = append(f.Elements, FormElement{labelText: "todo", elementType: FormElementType(eType)})
	} else {
		log.Printf("Form contains element I cannot work with yet: %s", n.Data)
	}
}
