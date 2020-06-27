package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func uniq(r io.Reader, w io.Writer) error {
	in := bufio.NewScanner(r)
	var prev string

	for in.Scan() {
		txt := in.Text()

		if txt == prev {
			continue
		}

		if txt < prev {
			return fmt.Errorf("file not sorted")
		}

		prev = txt

		fmt.Fprintln(w, txt)
	}

	return nil
}

func main() {
	err := uniq(os.Stdin, os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}
