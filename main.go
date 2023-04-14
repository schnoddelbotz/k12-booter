package main

import (
	"flag"
	"fmt"
	"os"
	"schnoddelbotz/k12-booter/cui"
	"schnoddelbotz/k12-booter/formgenerator"
	"schnoddelbotz/k12-booter/internationalization"
)

type flags struct {
	formByQuery  bool
	localeInfo   bool
	formFromFile string
}

func main() {
	fmt.Println("Welcome to nc-booter K12 EDU OSS IT Wizard. Please wait ... :)")

	options := &flags{}
	flag.BoolVar(&options.formByQuery, "formByQuery", false, "generate form code interactively")
	flag.StringVar(&options.formFromFile, "formFromFile", "", "generate form code from input file")
	flag.BoolVar(&options.localeInfo, "localeInfo", false, "dump detected locale info an quit job")
	flag.Parse()

	if options.localeInfo {
		info := internationalization.GetLocaleInfo()
		fmt.Printf("Detected locale: %s\n", info.Locale)
		fmt.Printf("Region (~country): %s, confidence: %s\n", info.Region, info.RegionConfidence)
		fmt.Printf("Detected base language: %s, confidence: %s\n", info.Base, info.BaseConfidence)
		fmt.Printf("Detected script: %s, confidence: %s\n", info.Script, info.ScriptConfidence)
		fmt.Printf("CultureInfo(756): %+v\n", internationalization.CultureInfo(756))
		fmt.Printf(`CultureInfo("Australia"): %+v`+"\n", internationalization.CultureInfo("Australia"))
		os.Exit(1)
	}

	if formgenerator.CreateFormAsNeeded(options.formByQuery, options.formFromFile) {
		os.Exit(0)
	}

	// Launch the UI
	cui.Zain()
}
