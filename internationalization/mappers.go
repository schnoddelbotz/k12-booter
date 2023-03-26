package internationalization

type InternalCountryCode int

const (
	ISOCountryUSA InternalCountryCode = iota
	ISOCountryGermany
	ISOCountrySwitzerland
	ISOCountryJapan
	ISOCountryBahamas
	ISOCountryGhana
	ISOCountryUAE
	ISOCountryChina
	ISOCountryRussia
	ISOCountryMexico
	ISOCountryPeru
	ISOCountryItaly
	ISOCountryTODO
)
const (
	LookupByFlag int = iota
	LookupByCountryName
	LookupByLanguage
	LookupByLocale
	LookupByShorthand
	LookupByInternalCountryCode
)

type Language struct {
	EnglishName string
	LocalName   string
	BCP47       string
}
type Locale struct {
	Name     string
	Country  Country
	Language Language
}
type Country struct {
	InternalCC      InternalCountryCode
	ISOName         string
	CommonShorthand string
	Languages       []Language
	LocalNames      []string
	Flag            string
}

var CountryCodeISOMAP = map[InternalCountryCode]string{
	ISOCountryChina: "China",
}

func AvailableTranslations() []Language {
	languages := []Language{}
	languages = append(languages, Language{BCP47: "en"})
	// deduce by reading embedded go files from translations
	return languages
}

func GetLocale(locale string) {
	// use UNIX locale, not http header.
	// auto detects from env, falls back to en always, as safest bet
}

func CheckConsistency() error {
	// ensure all maps work in all directions.
	// flag :germany: -> Germany
	// Germany -> flag :germany:
	// de -> Germany
	// Germany -> de
	return nil
}

// Stadt Konstanz funfact. We have a restaurant Volapük in Litzelstetten.
// See also: https://en.wikipedia.org/wiki/Codes_for_constructed_languages
// https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
// https://pkg.go.dev/golang.org/x/text/language
// https://en.wikipedia.org/wiki/IETF_language_tag
// An IETF BCP 47 language tag is a standardized code or tag that is used to identify human languages in the Internet.
// See also Hildegard von Bingen via Harry. Nature medicine applied.
