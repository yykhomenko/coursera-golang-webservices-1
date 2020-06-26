package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleHash(t *testing.T) {
	in := make(chan interface{}, 2)
	out := make(chan interface{}, 1)

	in <- 1
	in <- 2

	close(in)
	defer close(out)

	SingleHash(in, out)

	assert.Equal(t, "2212294583~709660146", <-out)
	assert.Equal(t, "450215437~1933333237", <-out)
}

func TestMultiHash(t *testing.T) {
	in := make(chan interface{}, 4)
	out := make(chan interface{}, 1)

	in <- "1"
	in <- "2"
	close(in)
	defer close(out)

	MultiHash(in, out)

	assert.Equal(t, "35962279594252452532383231384528719107062989936755", <-out)
	assert.Equal(t, "133085716516859850382103780943841265288725582281", <-out)
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
