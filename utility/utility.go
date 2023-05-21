package utility

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Fetch(url, file string) {
	if _, err := os.Stat(file); err == nil {
		// fmt.Printf("File %s exists, skipping download\n", file)
		return
	}

	out, err := os.Create(file)
	Fatal(err)
	defer out.Close()

	resp, err := http.Get(url)
	Fatal(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	Fatal(err)
	// fmt.Printf("Downloaded %s to %s (%d bytes) successfully\n", url, file, n)
}

func Fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
