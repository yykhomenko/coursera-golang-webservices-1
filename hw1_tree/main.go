package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

func (dir Dir) String() string {
	return dir.name
}

func (file File) String() string {
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
		if !(file.IsDir() || includeFiles) {
			continue
		}

		var newNode Node
		if file.IsDir() {
			children, err := readNodes(filepath.Join(path, file.Name()), []Node{}, includeFiles)
			if err != nil {
				continue
			}
			newNode = Dir{file.Name(), children}
		} else {
			newNode = File{file.Name(), file.Size()}
		}

		nodes = append(nodes, newNode)
	}

	return nodes, err
}

func printDir(out io.Writer, nodes []Node, prefixes []string) {
	if len(nodes) == 0 {
		return
	}

	fmt.Fprintf(out, "%s", strings.Join(prefixes, ""))

	node := nodes[0]

	if len(nodes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", node)
		if dir, ok := node.(Dir); ok {
			printDir(out, dir.children, append(prefixes, "\t"))
		}
		return
	}

	fmt.Fprintf(out, "%s%s\n", "├───", node)
	if dir, ok := node.(Dir); ok {
		printDir(out, dir.children, append(prefixes, "│\t"))
	}

	printDir(out, nodes[1:], prefixes)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
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
