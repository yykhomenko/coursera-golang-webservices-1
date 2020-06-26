package main

import (
	"log"
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
	}
	sort.Strings(data)
	result := strings.Join(data, "_")
	log.Println("combine results:", result)
	out <- result
}
