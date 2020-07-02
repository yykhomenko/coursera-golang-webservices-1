package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Name     string
	Email    string
	Browsers []string
}

var replacer = strings.NewReplacer("@", " [at] ")

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	// var seenBrowsers []string
	seenBrowsers := make(map[string]struct{})
	var userCounter int

	fmt.Fprintln(out, "found users:")

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		user := &User{}
		if err := json.Unmarshal(sc.Bytes(), user); err != nil {
			panic(err)
		}

		isMSIE := false
		isAndroid := false
		for _, browser := range user.Browsers {

			if strings.Contains(browser, "Android") {
				isAndroid = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
				}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
				}
			}
		}

		if isAndroid && isMSIE {
			email := replacer.Replace(user.Email)
			io.WriteString(out, fmt.Sprintf("[%d] %s <%s>\n", userCounter, user.Name, email))
		}

		userCounter++
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
