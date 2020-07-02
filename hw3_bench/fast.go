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
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		user := parseUser(sc.Bytes())
		users = append(users, user)
	}

	process(out, users)
}

func process(out io.Writer, users []User) {

	foundUsers := ""
	var seenBrowsers []string

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
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:")
	io.WriteString(out, foundUsers)
	io.WriteString(out, "\n")
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

func parseUser(line []byte) User {
	user := User{}
	if err := json.Unmarshal(line, &user); err != nil {
		panic(err)
	}
	return user
}
