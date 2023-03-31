package formgenerator

import "fmt"

func QueryFormParamsFromUser(doit bool) *Form {
	// dialog-driven, add elems, ask for filename to write to or stdout?
	if !doit {
		return nil
	}
	form := &Form{Complete: false}
	fElement := &FormElement{}
	var input string
	for !form.Complete {
		input = ""
		fmt.Print("Enter text: ")
		fmt.Scanln(&input)
		if input == "" {
			form.Complete = true
		} else {
			fElement.labelText = input
			fElement.elementType = FormElementType(FT_Input)
			form.Elements = append(form.Elements, *fElement)
		}
	}
	return form
}

func CreateFormAsNeeded(byQuery bool, byFilename string) bool {
	form := &Form{}
	if byQuery {
		form = QueryFormParamsFromUser(byQuery)
	}
	if byFilename != "" {
		form = ReadFormFromHTMLFile(byFilename)
	}
	if form.Complete {
		// HTML for now to debug import into *Form, Go later
		// This should reflect the original input form structure. Test...
		CreateFormHTML(form)
		return true
	}
	return false
}
