package aptbuddy

// todo
// index debian package infos (long description, tags (!)) for k12booter search
// --
// diff fsf osd vs debian packages ?

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

func OpenIndex(translation string) (a bleve.Index, e error) {
	a, e = bleve.Open(getIndexFilename(translation))
	return
}

func IndexExists(translation string) (bool, error) {
	idx, e := bleve.Open(getIndexFilename(translation))
	if e != nil {
		return false, e
	}
	return true, idx.Close()
}

func CreateIndex(translation string) error {
	mapping := bleve.NewIndexMapping()
	idx, err := bleve.New(getIndexFilename(translation), mapping)
	if err != nil {
		return err
	}
	return idx.Close()
}

func getIndexFilename(translation string) string {
	return fmt.Sprintf("aptbuddy_%s.bleve", translation) // path?
}
