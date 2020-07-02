package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/mailru/easyjson/jlexer"
)

var (
	cAndroid  = "Android"
	cMSIE     = "MSIE"
	cbAndroid = []byte(cAndroid)
	cbMSIE    = []byte(cMSIE)

	replacer = strings.NewReplacer("@", " [at] ")
	userPool = sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	}
)

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
	defer file.Close()

	fmt.Fprintln(out, "found users:")

	seenBrowsers := make(map[string]struct{}, 200)
	sc := bufio.NewScanner(file)
	for i := 0; sc.Scan(); i++ {

		line := sc.Bytes()

		if !(bytes.Contains(line, cbAndroid) || bytes.Contains(line, cbMSIE)) {
			continue
		}

		// user := &User{}

		user := userPool.Get().(*User)

		if err := user.UnmarshalJSON(line); err != nil {
			panic(err)
		}

		isMSIE := false
		isAndroid := false
		for _, browser := range user.Browsers {

			if strings.Contains(browser, cAndroid) {
				isAndroid = true
				seenBrowsers[browser] = struct{}{}
			}

			if strings.Contains(browser, cMSIE) {
				isMSIE = true
				seenBrowsers[browser] = struct{}{}
			}
		}

		if isAndroid && isMSIE {
			email := replacer.Replace(user.Email)
			io.WriteString(out, fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))
		}

		userPool.Put(user)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
