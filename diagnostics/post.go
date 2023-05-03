package diagnostics

import (
	"fmt"
	"io"
	"time"
)

// POST MAIN
func RunPOST(output io.Writer) {
	// run all tests
	// provide remedy instructions on fail
	// write all output to io.buf ( -> View.Main )
	for i := 0; i < 10; i++ {
		fmt.Fprintf(output, "POST %d", i)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Fprintf(output, "POST DONE")
}
