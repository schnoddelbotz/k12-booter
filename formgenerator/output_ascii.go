package formgenerator

import (
	"fmt"
	"log"
)

func CreateFormPlainText(f *Form) {
	log.Print("Creating Plaintext pseudo-lynx output from in-memory Form representation")
	fmt.Println("# FORM")
	for _, v := range f.Elements {
		fmt.Println(v.RenderPlainText(f))
	}
	fmt.Println("# END OF FORM")
}

// todo lynx-like html rendering of form, eases debugging hopefully ...
func (fe *FormElement) RenderPlainText(f *Form) string {
	// use real templates ...?!
	switch fe.elementType {
	case FT_Input:
		switch fe.inputType {
		case InputType_Text:
			return fmt.Sprintf(`%s____________ `, fe.labelText)
		case InputType_Submit:
			return `[ SUBMIT ]`
		case InputType_Checkbox:
			re := ""
			for _, c := range fe.selectOptions {
				label := labelForIdPlainText(f, c.id)
				checked := " "
				if label == "" {
					label = c.value
				}
				if c.selected {
					checked = "x"
				}
				re += fmt.Sprintf("[%s] %s ", checked, label)
			}
			return re
		case InputType_Radio:
			re := ""
			for _, c := range fe.selectOptions {
				label := labelForIdPlainText(f, c.id)
				checked := " "
				if label == "" {
					label = c.value
				}
				if c.selected {
					checked = "x"
				}
				re += fmt.Sprintf("(%s) %s ", checked, label)
			}
			return re
		default:
			return fmt.Sprintf(`%s_ _ _ _ _ _ _ `, fe.labelText)
		}
	case FT_Select:
		txt := fe.selectOptions[0].label
		for k, opt := range fe.selectOptions {
			if opt.selected {
				txt = fe.selectOptions[k].label
			}
		}
		return fmt.Sprintf(`%s [v_____%s]`, fe.labelText, txt)
	case FT_LineBreak:
		return ""
	case FT_TextElement:
		return fe.text
	default:
		return "??? BUG ???"
	}
}
