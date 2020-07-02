package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mailru/easyjson/jlexer"
)

var replacer = strings.NewReplacer("@", " [at] ")

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

func (out *User) UnmarshalJSON(data []byte) error {
	in := jlexer.Lexer{Data: data}
	unmarshal(&in, out)
	return in.Error()
}

func unmarshal(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	seenBrowsers := make(map[string]struct{})
	var userCounter int

	fmt.Fprintln(out, "found users:")

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		user := &User{}
		if err := user.UnmarshalJSON(sc.Bytes()); err != nil {
			panic(err)
		}

		isMSIE := false
		isAndroid := false
		for _, browser := range user.Browsers {

			if strings.Contains(browser, "Android") {
				isAndroid = true
				seenBrowsers[browser] = struct{}{}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				seenBrowsers[browser] = struct{}{}
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
