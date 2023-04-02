package formgenerator

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Nota bene - https://pkg.go.dev/golang.org/x/net/html -
// "It is the caller's responsibility to ensure that the Reader provides UTF-8 encoded HTML."

func ReadFormFromHTMLFile(filename string) *Form {
	data, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}
	defer data.Close()
	contents, err := io.ReadAll(data)
	if err != nil {
		log.Panic(err)
	}
	return ReadFormFromString(string(contents))
}

// ReadFormFromString translates an input (HTML+) <FORM> into internal *Form representation
func ReadFormFromString(htmlInput string) *Form {
	f := NewForm()
	htmlInputReader := strings.NewReader(htmlInput)
	doc, err := html.Parse(htmlInputReader)
	if err != nil {
		log.Panicf("fixme: don't panic here. html parser error. return err: %v", err)
	}

	// ~ 1:1 copy https://pkg.go.dev/golang.org/x/net/html example
	// formNode is "the" (...) <form> needle in provided HTML haystack, if found.
	var formNode *html.Node
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" {
			log.Printf("Found form. Form Attributes: %+v", n.Attr)
			log.Print(renderNode(n))
			formNode = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}
	fn(doc)

	// process all elements inside the <FORM> tag, add them to Form.Elements[]
	for c := formNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			nodeText := trimNewlinesAndWhitespace(renderNode(c))
			if nodeText != "" {
				f.AddTextAsElement(nodeText)
			}
			continue
		}
		if c.Type != html.ElementNode {
			continue
		}
		// log.Printf("Element: %s", renderNode(c))
		f.AddHTMLNodeAsElement(c)
	}
	f.SetLabelTextPerID()
	f.Complete = true
	return f
}

func renderNode(n *html.Node) string {
	// from: https://zetcode.com/golang/net-html/
	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := html.Render(w, n)
	if err != nil {
		return ""
	}
	return buf.String()
}
