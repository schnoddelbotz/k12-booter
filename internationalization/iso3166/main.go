package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type iso3166 struct {
	CountryName       string
	OfficialStateName string
	Sovereignty       string
	Alpha2Code        string
	Alpha3Code        string
	NumericCode       string
	SubdivisionCode   string
	InternetccTLD     string
}

type cliflags struct {
	in  string
	to  string
	via string
}

const outputFormatGO = "go"

func main() {
	flags := cliflags{}
	flag.StringVar(&flags.in, "in", "List_of_ISO_3166_country_codes.html", "local filename")
	flag.StringVar(&flags.to, "to", "go", "output format")
	flag.StringVar(&flags.via, "via", "https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes", "source url")

	flag.Parse()
	fetch(flags.via, flags.in)
	isodata := readhtml(flags.in)
	switch flags.to {
	case outputFormatGO:
		produceGo(isodata)
	default:
		fmt.Println("Ja, well. Only go is supported as output format right now.")
	}
}

func produceGo(data []iso3166) {

	fmt.Printf("var ISO2CountryName = map[string]string{\n")
	for _, country := range data {
		fmt.Printf(`  "%s": "%s",`+"\n", country.Alpha2Code, country.CountryName)
	}
	fmt.Printf("}\n\n")

	fmt.Printf("var ISO2OfficialStateName = map[string]string{\n")
	for _, country := range data {
		fmt.Printf(`  "%s": "%s",`+"\n", country.Alpha2Code, country.OfficialStateName)
	}
	fmt.Printf("}\n\n")

	fmt.Printf("var TLD2ISO2 = map[string]string{\n")
	for _, country := range data {
		fmt.Printf(`  "%s": "%s",`+"\n", country.InternetccTLD, country.Alpha2Code)
	}
	fmt.Printf("}\n\n")

	fmt.Printf("var ISO3Name = map[string]string{\n")
	for _, country := range data {
		fmt.Printf(`  "%s": "%s",`+"\n", country.Alpha3Code, country.Alpha2Code)
	}
	fmt.Printf("}\n\n")

	fmt.Printf("var ISO2Numeric = map[string]string{\n")
	for _, country := range data {
		fmt.Printf(`  "%s": "%s",`+"\n", country.Alpha2Code, country.NumericCode)
	}
	fmt.Printf("}\n\n")

	fmt.Printf("var Country2Sovereignity = map[string]string{\n")
	for _, country := range data {
		fmt.Printf(`  "%s": "%s",`+"\n", country.CountryName, country.Sovereignty)
	}
	fmt.Printf("}\n\n")

	// ~~~~~~~
	// dumpTableToStringMap(1,3) // i.e. provide column numbers, construct map name based on column name?
	// dumpTableToMap(1,[2,4,5,6]) // i.e. create a struct with "correct" automatic names
}

func readhtml(file string) []iso3166 {
	htmlFile, err := os.Open(file)
	fatal(err)
	defer htmlFile.Close()

	htmldoc, err := html.Parse(htmlFile)
	fatal(err)
	var tableNode *html.Node
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" &&
			strings.Contains(getElementAttributeValue(n.Attr, "class"), "wikitable") {
			tableNode = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}
	fn(htmldoc)
	if tableNode == nil {
		fatal(errors.New("no table of CSS class wikitable found in HTML; tried hard, very"))
	}
	return readtable(tableNode)
}

func readtable(table *html.Node) []iso3166 {
	result := []iso3166{}
	// log.Printf("F node = %+v", table)

	var tbody *html.Node
	for c := table.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "tbody" {
			tbody = c
		}
	}
	if tbody == nil || tbody.FirstChild == nil {
		fatal(errors.New("tbody not found"))
	}

	for d := tbody.FirstChild; d != nil; d = d.NextSibling {
		if d.Data == "tr" {
			r := readrow(d)
			if r != nil {
				result = append(result, *r)
			}
		}
	}

	return result
}

func readrow(tablerow *html.Node) *iso3166 {
	// log.Printf("table row %+v", tablerow)
	result := &iso3166{}
	tdCounter := 0
	for td := tablerow.FirstChild; td != nil; td = td.NextSibling {
		if td.Type != html.ElementNode || td.Data == "th" {
			continue
		}
		switch tdCounter {
		case 0:
			result.CountryName = renderNode(deepest(td.LastChild.PrevSibling.FirstChild))
		case 1:
			result.OfficialStateName = renderNode(deepest(td))
		case 2:
			result.Sovereignty = strings.TrimSuffix(renderNode(deepest(td)), "\n")
		case 3:
			result.Alpha2Code = renderNode(deepest(td.FirstChild.LastChild))
		case 4:
			result.Alpha3Code = renderNode(td.FirstChild.LastChild.FirstChild)
		case 5:
			result.NumericCode = renderNode(td.FirstChild.LastChild.FirstChild)
		case 6:
			result.SubdivisionCode = renderNode(deepest(td))
		case 7:
			result.InternetccTLD = renderNode(deepest(td))
		}
		tdCounter++
	}
	if tdCounter != 8 {
		return nil
	}
	return result
}

func deepest(node *html.Node) *html.Node {
	for node.FirstChild != nil {
		return deepest(node.FirstChild)
	}
	return node
}

func fetch(url, file string) {
	if _, err := os.Stat(file); err == nil {
		// fmt.Printf("File %s exists, skipping download\n", file)
		return
	}

	out, err := os.Create(file)
	fatal(err)
	defer out.Close()

	resp, err := http.Get(url)
	fatal(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	fatal(err)
	// fmt.Printf("Downloaded %s to %s (%d bytes) successfully\n", url, file, n)
}

func fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// copy from import_html.go :/
func renderNode(n *html.Node) string {
	if n == nil {
		return ""
	}
	// from: https://zetcode.com/golang/net-html/
	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := html.Render(w, n)
	if err != nil {
		return ""
	}
	return buf.String()
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
