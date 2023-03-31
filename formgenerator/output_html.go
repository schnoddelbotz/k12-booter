package formgenerator

import (
	"fmt"
	"log"
)

func CreateFormHTML(f *Form) {
	log.Printf("PRODUCING HTML for Form ... %+v", f)
	// can be used to diff against imported html form ...
	fmt.Println("<form>")
	for _, v := range f.Elements {
		fmt.Println(v.RenderHTML())
	}
	fmt.Println("</form>")
}

func (fe *FormElement) RenderHTML() string {
	switch fe.elementType {
	case FT_Input:
		return fmt.Sprintf(`<input type="%s" name="%s" />`, fe.GetInputTypeName(), fe.name)
	case FT_Select:
		return fmt.Sprintf(`<select name="%s" />  ..tbd... </select>`, fe.name)
	default:
		return "HTML:N/A=FIXME"
	}
}
