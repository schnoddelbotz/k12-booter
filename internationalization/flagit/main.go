package main

import (
	"bufio"
	"fmt"
	"os"
	"schnoddelbotz/k12-booter/internationalization"
)

func main() {
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
				flag := flags[flagKey]
				flagKey++
				if flagKey == numFlags {
					flagKey = 0
				}
				fmt.Printf("%s", flag)
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
