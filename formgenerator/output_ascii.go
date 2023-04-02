package formgenerator

import (
	"fmt"
	"log"
)

func CreateFormPlainText(f *Form) {
	log.Print("Creating Plaintext pseudo-lynx output from in-memory Form representation")
	fmt.Println("# FORM")
	for _, v := range f.Elements {
		fmt.Println(v.RenderPlainText())
	}
	fmt.Println("# END OF FORM")
}

// todo lynx-like html rendering of form, eases debugging hopefully ...
func (fe *FormElement) RenderPlainText() string {
	// use real templates ...?!
	switch fe.elementType {
	case FT_Input:
		if fe.inputType == InputType_Text {
			return fmt.Sprintf(`%s____________ `, fe.labelText)
		} else if fe.inputType == InputType_Submit {
			return `[ SUBMIT ]`
		} else {
			return fmt.Sprintf(`%s_ _ _ _ _ _ _ `, fe.labelText)
		}
	case FT_Select:
		txt := fe.selectOptions[0].label
		for k, opt := range fe.selectOptions {
			if opt.selected {
				txt = fe.selectOptions[k].label
			}
		}
		return fmt.Sprintf(`%s[_____%s]`, fe.labelText, txt)
	case FT_LineBreak:
		return ""
	case FT_TextElement:
		return fe.text
	default:
		return "??? BUG ???"
	}
}
