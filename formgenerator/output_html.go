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

func (fe *FormElement) RenderHTML() string {
	// use real templates ...?!
	switch fe.elementType {
	case FT_Input:
		return fmt.Sprintf(`<input type="%s" name="%s" />`, fe.GetInputTypeName(), fe.name)
	case FT_Select:
		tpl := fmt.Sprintf(`<select name="%s">`+"\n", fe.name)
		for _, opt := range fe.selectOptions {
			selected := ""
			if opt.selected {
				selected = " selected"
			}
			tpl += fmt.Sprintf(`  <option value="%s"%s>%s</option>`+"\n", opt.value, selected, opt.label)
		}
		tpl += `</select>`
		return tpl
	default:
		return "HTML:N/A=FIXME"
	}
}
