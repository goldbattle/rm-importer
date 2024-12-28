package network

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

/*
Downloads all urls passed as arguments and passes them directly to 1Password CLI
*/

func fetch0() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch error %v\n", err)
		}

		cmd := exec.Command("op", "document", "create", "--title", "hello world", "--file-name", "hello world.txt", "--format", "json")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer stdin.Close()
			io.Copy(stdin, resp.Body)
		}()

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", out)
	}
}
