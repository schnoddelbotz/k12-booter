package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
	fmt.Printf("type var CountryFlags = map[string]string{ ...\n")
	for _, country := range data {
		fmt.Printf("c:%s -> tld:%s\n", country.CountryName, country.InternetccTLD)
	}
	fmt.Printf("} todo ... and all variants we need.")
	// dumpTableToStringMap(1,3) // i.e. provide column numbers, construct map name based on column name?
	// dumpTableToMap(1,[2,4,5,6]) // i.e. create a struct with "correct" automatic names
}

func readhtml(file string) []iso3166 {
	result := []iso3166{}

	htmlFile, err := os.Open(file)
	fatal(err)
	defer htmlFile.Close()

	htmldoc, err := html.Parse(htmlFile)
	fatal(err)
	var tableNode *html.Node
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			log.Printf("Found table. Attributes: %+v", n.Attr)
			log.Print(renderNode(n))
			tableNode = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}
	fn(htmldoc)

	log.Printf("F node = %+v", tableNode)

	result = append(result, iso3166{CountryName: "Germany", InternetccTLD: "de"})

	return result
}

func fetch(url, file string) {
	if _, err := os.Stat(file); err == nil {
		fmt.Printf("File %s exists, skipping download\n", file)
		return
	}

	out, err := os.Create(file)
	fatal(err)
	defer out.Close()

	resp, err := http.Get(url)
	fatal(err)
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	fatal(err)
	fmt.Printf("Downloaded %s to %s (%d bytes) successfully\n", url, file, n)
}

func fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// copy from import_html.go :/
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
