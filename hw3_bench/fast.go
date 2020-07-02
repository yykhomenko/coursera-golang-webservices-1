package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type User struct {
	Name     string
	Email    string
	Browsers []string
}

var atRegex = regexp.MustCompile("@")

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	var users []User
	var seenBrowsers []string

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		user := parseUser(sc.Bytes())
		users = append(users, user)
	}

	fmt.Fprintln(out, "found users:")

	var idx int
	for i, user := range users {

		isAndroid := false
		isMSIE := false
		for _, browser := range user.Browsers {

			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := atRegex.ReplaceAllString(user.Email, " [at] ")
		io.WriteString(out, fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

func parseUser(line []byte) User {
	user := User{}
	if err := json.Unmarshal(line, &user); err != nil {
		panic(err)
	}
	return user
}
