package main

import (
	"flag"
	"fmt"
	"os"
	"schnoddelbotz/k12-booter/aptbuddy"
	"schnoddelbotz/k12-booter/cui"
	"schnoddelbotz/k12-booter/formgenerator"
	"schnoddelbotz/k12-booter/internationalization"
	"schnoddelbotz/k12-booter/utility"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type flags struct {
	formByQuery           bool
	localeInfo            bool
	runAptBleveExperiment bool
	formFromFile          string
}

var AppVersion = "git-0.0.0"

func main() {
	viper.SetEnvPrefix("K12B")
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".k12booter")
	viper.ReadInConfig()

	options := &flags{}
	flag.BoolVar(&options.formByQuery, "formByQuery", false, "generate form code interactively")
	flag.StringVar(&options.formFromFile, "formFromFile", "", "generate form code from input file")
	flag.BoolVar(&options.localeInfo, "localeInfo", false, "dump detected locale info an quit job")
	flag.BoolVar(&options.runAptBleveExperiment, "apt", false, "run apt bleve indexing experiment")
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

	if options.runAptBleveExperiment {
		b, err := aptbuddy.New("en")
		utility.Fatal(err)
		b.Experiments()
		os.Exit(0)
	}

	if formgenerator.CreateFormAsNeeded(options.formByQuery, options.formFromFile) {
		os.Exit(0)
	}

	// Launch the UI
	cui.Zain(AppVersion)
}
