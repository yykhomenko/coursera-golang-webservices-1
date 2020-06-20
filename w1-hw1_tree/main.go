package main

import (
	"fmt"
	"os"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out *os.File, path string, printFiles bool) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	names, err := f.Readdirnames(0)
	if err != nil {
		return err
	}

	fmt.Fprintln(out, names)

	return nil
}
