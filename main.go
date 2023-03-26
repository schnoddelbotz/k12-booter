package main

import (
	"fmt"
	"time"

	"schnoddelbotz/k12-booter/cui"
)

func main() {
	fmt.Println("Welcome to nc-booter K12 EDU OSS IT Wizard. Please wait ... :)")
	time.Sleep(1 * time.Second)

	// Launch the UI
	cui.Zain()
}
