package main

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSingleHash(t *testing.T) {
	in := make(chan interface{})
	out := make(chan interface{})
	defer close(in)
	defer close(out)

	SingleHash(in, out)

	start := time.Now()

	in <- 1
	in <- 2

	assert.Equal(t, "2212294583~709660146", <-out)
	assert.Equal(t, "450215437~1933333237", <-out)

	log.Println("time:", time.Now().Sub(start))
}

func TestMultiHash(t *testing.T) {
	in := make(chan interface{})
	out := make(chan interface{})
	defer close(in)
	defer close(out)

	MultiHash(in, out)

	start := time.Now()

	in <- "1"
	in <- "2"

	assert.Equal(t, "347715282235962279594252452532383231384528719107062989936755", <-out)
	assert.Equal(t, "1447589260133085716516859850382103780943841265288725582281", <-out)

	log.Println("time:", time.Now().Sub(start))
}

func TestCombineResults(t *testing.T) {
	in := make(chan interface{}, 4)
	out := make(chan interface{}, 1)

	in <- "1"
	in <- "3"
	in <- "4"
	in <- "2"
	close(in)
	defer close(out)

	CombineResults(in, out)

	assert.Equal(t, "1_2_3_4", <-out)
}

func TestExecutePipeline(t *testing.T) {
	inputData := []int{0}
	testExpected := "1173136728138862632818075107442090076184424490584241521304_1696913515191343735512658979631549563179965036907783101867_27225454331033649287118297354036464389062965355426795162684_29568666068035183841425683795340791879727309630931025356555_3994492081516972096677631278379039212655368881548151736_4958044192186797981418233587017209679042592862002427381542_4958044192186797981418233587017209679042592862002427381542"
	testResult := "NOT_SET"

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				t.Error("cant convert result data to string")
			}
			testResult = data
		}),
	}

	// start := time.Now()

	ExecutePipeline(hashSignJobs...)

	// end := time.Since(start)
	//
	// expectedTime := 3 * time.Second

	if testExpected != testResult {
		t.Errorf("results not match\nGot: %v\nExpected: %v", testResult, testExpected)
	}

	// if end > expectedTime {
	// 	t.Errorf("execition too long\nGot: %s\nExpected: <%s", end, time.Second*3)
	// }
}
