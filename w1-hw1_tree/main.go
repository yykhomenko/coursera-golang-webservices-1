package main

import (
	"io"
	"os"
	"sort"
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

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for i, file := range files {

		prefix := "├───"
		if i == len(files)-1 {
			prefix = "└───"
		}

		if file.IsDir() {
			io.WriteString(out, prefix+file.Name())
			io.WriteString(out, "\n")

			printDir(out, path+string(os.PathSeparator)+file.Name(), printFiles)
		}
	}

	return nil
}
