package formgenerator

import (
	"fmt"
	"log"
)

func CreateFormHTML(f *Form) {
	log.Printf("Creating HTML <FORM> from in-memory Form representation: %+v", f)
	fmt.Println("<form>")
	for _, v := range f.Elements {
		fmt.Println(v.RenderHTML())
	}
	fmt.Println("</form>")
}

func labelFor(fe *FormElement) string {
	if fe.labelText == "" {
		return ""
	}
	return fmt.Sprintf(`<label for="%s">%s</label>`+"\n  ", fe.id, fe.labelText)
}

func (fe *FormElement) RenderHTML() string {
	// use real templates ...?!
	switch fe.elementType {
	case FT_Input:
		return fmt.Sprintf(`  %s<input id="%s" type="%s" name="%s" />`,
			labelFor(fe), fe.id, fe.GetInputTypeName(), fe.name)
	case FT_Select:
		tpl := fmt.Sprintf(`  %s<select name="%s" id="%s">`+"\n", labelFor(fe), fe.name, fe.id)
		for _, opt := range fe.selectOptions {
			selected := ""
			if opt.selected {
				selected = " selected"
			}
			tpl += fmt.Sprintf(`    <option value="%s"%s>%s</option>`+"\n", opt.value, selected, opt.label)
		}
		tpl += `  </select>`
		return tpl
	case FT_LineBreak:
		return "  <br>"
	case FT_Label:
		return `  ` + fe.text
	case FT_TextElement:
		return `  ` + fe.text
	default:
		return "HTML:N/A=FIXME"
	}
}
