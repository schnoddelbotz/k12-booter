package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"schnoddelbotz/k12-booter/internationalization"
)

var flagKey = 0
var flags []string
var numFlags int

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	re := regexp.MustCompile(`[^\s]`)
	flags = flagsAsArray()
	numFlags = len(flags)
	for scanner.Scan() {
		in := scanner.Text()
		in = re.ReplaceAllStringFunc(in, getFlag)
		fmt.Println(in)
	}

	if scanner.Err() != nil {
		panic(scanner.Err().Error())
	}
}

func getFlag(in string) string {
	flag := flags[flagKey]
	flagKey++
	if flagKey == numFlags {
		flagKey = 0
	}
	return flag
}

func flagsAsArray() (flags []string) {
	for _, c := range internationalization.Cultures {
		flags = append(flags, c.Flag)
	}
	return
}
