package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
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
		produceGo1(isodata)
	default:
		fmt.Println("Ja, well. Only go is supported as output format right now.")
	}
}

func produceGo1(data []iso3166) {
	// https://www.c-sharpcorner.com/UploadFile/0c1bb2/display-country-list-without-database-in-Asp-Net-C-Sharp/

	fmt.Println(`
package internationalization

// this file is auto-generated. do not hand-edit. Use Makefile to re-generate.

type CountryData struct {
	CountryName       string
	OfficialStateName string
	Sovereignty       string
	Alpha2Code        string
	Alpha3Code        string
	NumericCode       int
	SubdivisionCode   string
	InternetccTLD     string 
	Flag string
}`)
	fmt.Printf("var Cultures = []CountryData{\n")
	for _, country := range data {
		numCode, err := strconv.ParseUint(country.NumericCode, 10, 16)
		fatal(err)
		fmt.Printf(`{  
			Alpha2Code: "%s",
			Alpha3Code: "%s",
			CountryName: "%s",
			OfficialStateName: "%s",
			NumericCode: %d,
			InternetccTLD: "%s",
			Flag: Flag_%s,
		},
		`, country.Alpha2Code,
			country.Alpha3Code,
			country.CountryName,
			//strings.Replace(country.OfficialStateName, "&#39;", "'", 1),
			country.OfficialStateName,
			numCode, // country.NumericCode,
			country.InternetccTLD,
			CountryNameToFlagConstant(country.CountryName))
	}
	fmt.Printf("\n}\n")
}

func CountryNameToFlagConstant(cn string) string {
	// no politics intended here - just to match our Go constant names for FLAGS #PEACE.
	cn = strings.Replace(cn, ", the United Republic of", "", 1)
	cn = strings.Replace(cn, ", Kingdom of", "", 1)
	cn = strings.Replace(cn, "(", "", 2)
	cn = strings.Replace(cn, ")", "", 2)
	cn = strings.Replace(cn, " and ", "_", 1)
	cn = strings.Replace(cn, " the", "", 2)
	cn = strings.Replace(cn, " Federated States of", "", 1)
	cn = strings.Replace(cn, " Republic of", "", 1)
	cn = strings.Replace(cn, ", State of", "", 1)
	//cn = strings.Replace(cn, " State of", "", 1)
	cn = regexp.MustCompile(`\([^)]+\)`).ReplaceAllString(cn, "")
	cn = strings.Replace(cn, "-", "_", 1)
	cn = strings.Replace(cn, " ", "_", 9)
	cn = strings.Replace(cn, "'", "", 1) // e.g. Cote d'Ivoire
	cn = strings.Replace(cn, "Saint_", "St", 1)
	cn = strings.Replace(cn, ".", "", 2)           // U.S.
	cn = strings.Replace(cn, "_[Malvinas]", "", 1) // Falkland Islands
	return cn
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
			if string(result.CountryName[0]) == "[" {
				result.CountryName = renderNode(td.LastChild.PrevSibling.PrevSibling.PrevSibling.FirstChild)
				result.CountryName = strings.Replace(result.CountryName, "&#39;", "'", 1)
			}

		case 1:
			result.OfficialStateName = renderNode(deepest(td))
			result.OfficialStateName = strings.Replace(result.OfficialStateName, "&#39;", "'", 1)
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
			if string(result.InternetccTLD[0]) == "[" {
				result.InternetccTLD = "???"
			}
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
