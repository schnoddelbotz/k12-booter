package aptbuddy

import (
	"fmt"
	"schnoddelbotz/k12-booter/utility"

	"github.com/blevesearch/bleve/v2"
)

var (
	HighlightNone        = []string{}
	HighlightPackage     = []string{"Package"}
	HighlightDescription = []string{"Description"}
	HighlightDescAndTags = []string{"Description", "Tags"}
)

type Buddy struct {
	index bleve.Index
}

func New(translation string) (*Buddy, error) {
	var bud Buddy
	var err error

	// re-use existing .bleve index
	if exists, _ := IndexExists(translation); exists {
		bud.index, err = OpenIndex(translation)
		fmt.Println("re-using existing .bleve index")
		return &bud, err
	}

	// if no .bleve existed, download from debian ...
	fmt.Printf("fetching packages+translation_%s.gz from debain ...\n", translation)
	err = FetchFiles(translation)
	if err != nil {
		return &bud, err
	}
	// ...  and index
	fmt.Println("building .bleve index ...")
	err = CreateIndex(translation)
	utility.Fatal(err)
	bud.index, err = OpenIndex(translation)
	utility.Fatal(err)
	bud.Debian2Bleve(translation)
	fmt.Printf("created new index, now open with %+v", bud.index.Stats())
	return &bud, err
}

func (buddy *Buddy) Search(q string, s int, hilight []string) *bleve.SearchResult {
	query := bleve.NewQueryStringQuery(q)
	searchRequest := bleve.NewSearchRequest(query)

	// https://blevesearch.com/docs/Highlight%20Matches%20in%20Results/
	if len(hilight) > 0 {
		searchRequest.Highlight = bleve.NewHighlight()
		for _, hl := range hilight {
			searchRequest.Highlight.AddField(hl)
		}
	}

	searchRequest.Size = s
	fmt.Printf("Search: `%s`, showing max %d results ...\n", q, s)
	searchResult, _ := buddy.index.Search(searchRequest)
	return searchResult
}

func (buddy *Buddy) FieldDict(fieldname string, minCount uint64, maxCount uint64) error {
	// = https://github.com/blevesearch/bleve/blob/master/cmd/bleve/cmd/dictionary.go
	fmt.Printf("Dumping dictionary for field %s, for terms with a minium count of %d, max %d\n",
		fieldname, minCount, maxCount)
	i, err := buddy.index.Advanced()
	if err != nil {
		return err
	}
	ar, err := i.Reader()
	if err != nil {
		return err
	}
	d, err := ar.FieldDict(fieldname)
	if err != nil {
		return err
	}
	de, err := d.Next()
	for err == nil && de != nil {
		if de.Count > minCount && de.Count < maxCount {
			fmt.Printf("%s - %d\n", de.Term, de.Count)
			// todo: don't print, return sorted
		}
		de, err = d.Next()
	}
	println()
	return nil
}

func (buddy *Buddy) Close() error {
	return buddy.index.Close()
}
