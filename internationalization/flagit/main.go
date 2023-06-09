package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"schnoddelbotz/k12-booter/internationalization"
)

func main() {
	userFlag := flag.String("flag", "", "Use given emoji/flag instead of all")
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	flags := flagsAsArray()
	numFlags := len(flags)
	flagKey := 0
	for scanner.Scan() {
		in := scanner.Text()
		for _, c := range in {
			if c == ' ' {
				fmt.Printf(" ")
			} else {
				if *userFlag == "" {
					flag := flags[flagKey]
					flagKey++
					if flagKey == numFlags {
						flagKey = 0
					}
					fmt.Printf("%s", flag)
				} else {
					fmt.Printf("%s", *userFlag)
				}
			}
		}
		fmt.Println()
	}

	if scanner.Err() != nil {
		panic(scanner.Err().Error())
	}
}

func flagsAsArray() (flags []string) {
	for _, c := range internationalization.Cultures {
		flags = append(flags, c.Flag)
	}
	return
}
