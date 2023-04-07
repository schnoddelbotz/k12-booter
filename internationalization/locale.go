package internationalization

import (
	"log"

	"github.com/Xuanwo/go-locale"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type LocaleInfo struct {
	Locale           string
	LanguageTag      language.Tag
	Region           language.Region
	RegionConfidence language.Confidence
	Base             language.Base
	BaseConfidence   language.Confidence
	Script           language.Script
	ScriptConfidence language.Confidence
	// TODO:
	CountryTwoLetter string
	CountryNameEN    string
	CountryNameLocal string
	LanguagesSpoken  []string
}

func GetLocaleInfo() LocaleInfo {
	lang, err := locale.Detect()
	if err != nil {
		log.Fatal(err)
	}
	info := LocaleInfo{Locale: lang.String(), LanguageTag: lang}
	info.Region, info.RegionConfidence = lang.Region()
	info.Base, info.BaseConfidence = lang.Base()
	info.Script, info.ScriptConfidence = lang.Script()

	log.Printf("%+v", info.LanguageTag)

	//x := display.Scripts(lang)
	log.Printf("XXX %+v", display.Scripts(lang)) // LATIN
	log.Printf("XXX %+v", display.German.Languages().Name(lang))
	// HOW TO GET ALL languages SPOKEN / available in e.g. Switzerland?

	var matcher = language.NewMatcher([]language.Tag{})
	tag, index, confidence := matcher.Match(info.LanguageTag)
	log.Printf("TAG from matcher %+v", tag)
	log.Printf("best match: %s (%s) index=%d confidence=%v\n",
		display.English.Tags().Name(info.LanguageTag),
		display.Self.Name(info.LanguageTag),
		index, confidence)

	//log.Printf("LLLL %+v", display.Languages())

	return info
}
