package main

import (
	"sort"
	"strconv"
	"strings"
)

const TH = 6

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})

	for _, job := range jobs {
		out := make(chan interface{})
		job(in, out)
		in = out
	}
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

			data := make([]string, TH)

			for i := 1; i < TH; i++ {
				crc32 := DataSignerCrc32(strconv.Itoa(i) + raw.(string))
				data = append(data, crc32)
			}

			out <- strings.Join(data, "")
		}
	}()
}

func CombineResults(in, out chan interface{}) {
	var data []string
	for raw := range in {
		data = append(data, raw.(string))
	}
	sort.Strings(data)
	out <- strings.Join(data, "_")
}
