package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
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
	m := &sync.Mutex{}
	go func() {
		for raw := range in {
			data := strconv.Itoa(raw.(int))
			go calcSingleHash(m, out, data)
		}
	}()
}

func calcSingleHash(m *sync.Mutex, out chan interface{}, data string) {
	md5ch := make(chan string)
	crc32ch := make(chan string)
	crc32md5ch := make(chan string)
	defer close(md5ch)
	defer close(crc32ch)
	defer close(crc32md5ch)

	go func(out chan<- string, data string) {
		m.Lock()
		out <- DataSignerMd5(data)
		m.Unlock()
	}(md5ch, data)

	go func(out chan<- string, data string) {
		out <- DataSignerCrc32(data)
	}(crc32ch, data)

	go func(out chan<- string, data string) {
		out <- DataSignerCrc32(data)
	}(crc32md5ch, <-md5ch)

	out <- (<-crc32ch) + "~" + (<-crc32md5ch)
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
