package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type Node interface {
	fmt.Stringer
}

type Dir struct {
	name     string
	children []Node
}

type File struct {
	name string
	size int64
}

func (dir *Dir) String() string {
	return dir.name
}

func (file *File) String() string {
	if file.size == 0 {
		return file.name + " (empty)"
	}

	return file.name + " (" + strconv.FormatInt(file.size, 10) + "b)"
}

func readNodes(path string, nodes []Node, includeFiles bool) ([]Node, error) {
	files, err := ioutil.ReadDir(path)

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {

	}

	return nodes, err
}

func printDir(out io.Writer, nodes []Node, prefixes []string) {

}

func dirsOnly(files []os.FileInfo) []os.FileInfo {
	b := make([]os.FileInfo, 0)
	for _, f := range files {
		if f.IsDir() {
			b = append(b, f)
		}
	}
	return b
}

func dirTree(out *os.File, path string, printFiles bool) error {
	nodes, err := readNodes(path, []Node{}, printFiles)
	printDir(out, nodes, []string{})

	return err
}

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
