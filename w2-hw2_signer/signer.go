package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

const TH = 6

func ExecutePipeline(jobs ...job) {
	var wg sync.WaitGroup
	in := make(chan interface{})
	out := make(chan interface{})

	wg.Add(len(jobs))
	for _, job := range jobs {
		go execJob(&wg, job, in, out)
		in = out
		out = make(chan interface{})
	}

	wg.Wait()
	close(out)
}

func execJob(wg *sync.WaitGroup, job job, in, out chan interface{}) {
	job(in, out)
	wg.Done()
	close(out)
}

func SingleHash(in, out chan interface{}) {
	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for raw := range in {
		wg.Add(1)
		go calcSingleHash(m, wg, out, strconv.Itoa(raw.(int)))
	}
	wg.Wait()
}

func calcSingleHash(m *sync.Mutex, wg *sync.WaitGroup, out chan interface{}, data string) {
	defer wg.Done()
	md5ch := make(chan string)
	crc32ch := make(chan string)
	crc32md5ch := make(chan string)

	go func(out chan<- string, data string) {
		m.Lock()
		out <- DataSignerMd5(data)
		m.Unlock()
		close(out)
	}(md5ch, data)

	go func(out chan<- string, data string) {
		out <- DataSignerCrc32(data)
		close(out)
	}(crc32ch, data)

	go func(out chan<- string, data string) {
		out <- DataSignerCrc32(data)
		close(out)
	}(crc32md5ch, <-md5ch)

	out <- (<-crc32ch) + "~" + (<-crc32md5ch)
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for raw := range in {
		wg.Add(1)
		go calcMultiHash(wg, out, raw.(string))
	}
	wg.Wait()
}

func calcMultiHash(wg *sync.WaitGroup, out chan interface{}, data string) {
	defer wg.Done()

	hashes := make([]chan string, TH)

	for i := 0; i < TH; i++ {
		hash := make(chan string)
		hashes[i] = hash
		go func(out chan string, idx int) {
			out <- DataSignerCrc32(strconv.Itoa(idx) + data)
			close(out)
		}(hash, i)
	}

	var buff []string
	for i := 0; i < TH; i++ {
		buff = append(buff, <-hashes[i])
	}

	out <- strings.Join(buff, "")
}

func CombineResults(in, out chan interface{}) {
	var data []string
	for raw := range in {
		data = append(data, raw.(string))
	}
	sort.Strings(data)
	out <- strings.Join(data, "_")
}
