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

type User map[string]interface{}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	users := make([]User, 0)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		user := parseUser(sc.Text())
		users = append(users, user)
	}

	foundUsers := ""
	var seenBrowsers []string
	rxp := regexp.MustCompile("@")

	for i, user := range users {

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			continue
		}

		isAndroid := false
		isMSIE := false
		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				continue
			}

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

		email := rxp.ReplaceAllString(user["email"].(string), " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
	}

	fmt.Fprintln(out, "found users:")
	io.WriteString(out, foundUsers)
	io.WriteString(out, "\n")
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

func parseUser(line string) User {
	user := make(User)
	if err := json.Unmarshal([]byte(line), &user); err != nil {
		panic(err)
	}
	return user
}
