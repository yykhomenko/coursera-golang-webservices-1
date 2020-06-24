package main

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const TH = 6

func ExecutePipeline(jobs ...job) {

}

func SingleHash(in, out chan interface{}) {
	go func() {
		for raw := range in {
			data := strconv.Itoa(raw.(int))
			out <- DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
		}
	}()
}

func MultiHash(in, out chan interface{}) {
	go func() {
		for raw := range in {
			data := raw.(string)
			b := strings.Builder{}

			for i := 1; i < TH; i++ {
				b.WriteString(DataSignerCrc32(strconv.Itoa(i) + data))
			}

			out <- b.String()
		}
	}()
}

func CombineResults(in, out chan interface{}) {
	var data []string
	for raw := range in {
		data = append(data, raw.(string))
		runtime.Gosched()
	}

	sort.Strings(data)
	// out <- strings.Join(data, "_")
	out <- "sd"
}

func main() {
	in := make(chan interface{}, 1)
	out := make(chan interface{}, 1)

	in <- "string"
	close(in)
	CombineResults(in, out)

	fmt.Println(<-out)
}
