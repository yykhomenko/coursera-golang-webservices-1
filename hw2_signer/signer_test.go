package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleHash(t *testing.T) {
	in := make(chan interface{})
	out := make(chan interface{})
	defer close(out)
	defer close(in)

	go SingleHash(in, out)

	in <- 1
	assert.Equal(t, "2212294583~709660146", <-out)
	in <- 2
	assert.Equal(t, "450215437~1933333237", <-out)
}

func TestMultiHash(t *testing.T) {
	in := make(chan interface{})
	out := make(chan interface{})
	defer close(out)
	defer close(in)

	go MultiHash(in, out)

	in <- "2212294583~709660146"
	assert.Equal(t, "4958044192186797981418233587017209679042592862002427381542", <-out)
	in <- "450215437~1933333237"
	assert.Equal(t, "27225454331033649287118297354036464389062965355426795162684", <-out)
}

func TestCombineResults(t *testing.T) {
	in := make(chan interface{}, 2)
	out := make(chan interface{})
	defer close(out)

	go CombineResults(in, out)

	in <- "4958044192186797981418233587017209679042592862002427381542"
	in <- "27225454331033649287118297354036464389062965355426795162684"
	close(in)

	assert.Equal(t, "27225454331033649287118297354036464389062965355426795162684_4958044192186797981418233587017209679042592862002427381542", <-out)
}
