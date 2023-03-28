package main

import (
	"flag"
	"fmt"
	"os"
	"schnoddelbotz/k12-booter/cui"
	"schnoddelbotz/k12-booter/formgenerator"
	"time"
)

type flags struct {
	formByQuery  bool
	formFromFile string
}

func main() {
	fmt.Println("Welcome to nc-booter K12 EDU OSS IT Wizard. Please wait ... :)")
	time.Sleep(1 * time.Second)

	options := &flags{}
	flag.BoolVar(&options.formByQuery, "formByQuery", false, "generate form code interactively")
	flag.StringVar(&options.formFromFile, "formFromFile", "", "generate form code from input file")
	flag.Parse()

	if formgenerator.CreateFormAsNeeded(options.formByQuery, options.formFromFile) {
		os.Exit(0)
	}

	// Launch the UI
	cui.Zain()
}
