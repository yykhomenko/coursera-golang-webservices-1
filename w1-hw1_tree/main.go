package main

import (
	"io"
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

	return printDir(out, path, printFiles)
}

func printDir(out *os.File, path string, printFiles bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	io.WriteString(out, f.Name())
	io.WriteString(out, "\n")

	for _, file := range files {
		if file.IsDir() {
			printDir(out, path+string(os.PathSeparator)+file.Name(), printFiles)
		}

		io.WriteString(out, file.Name())
		io.WriteString(out, "\n")
	}

	return nil
}
