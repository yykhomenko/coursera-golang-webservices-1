package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombineResults(t *testing.T) {
	in := make(chan interface{}, 4)
	out := make(chan interface{}, 1)

	in <- "1"
	in <- "3"
	in <- "4"
	in <- "2"
	close(in)

	CombineResults(in, out)

	assert.Equal(t, "1_2_3_4", <-out)
}
