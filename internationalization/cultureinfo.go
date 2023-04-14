package internationalization

import (
	"fmt"
)

// https://learn.microsoft.com/en-us/dotnet/api/system.globalization.cultureinfo.-ctor?view=net-8.0

func CultureInfo(x interface{}) *CountryData {
	// Convinience bloat? Provide SAME interface. Once. I thought.
	switch x := x.(type) {
	case int:
		return CultureInfoById(x)
	case string:
		return CultureInfoByName(x)
	default:
		panic(fmt.Sprintf("CultureInfo called with unexpected type %T", x))
	}
}

func CultureInfoById(id int) *CountryData {
	for _, v := range Cultures {
		if v.NumericCode == id {
			return &v
		}
	}
	return nil
}

func CultureInfoByName(name string) *CountryData {
	for _, v := range Cultures {
		if v.CountryName == name {
			return &v
		}
	}
	return nil
}
