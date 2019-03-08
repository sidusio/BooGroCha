package main

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestDaysToAdd(t *testing.T) {
	assert.Equal(t, daysToAdd(time.Saturday, time.Monday), 2)
	assert.Equal(t, daysToAdd(time.Monday, time.Saturday), 5)
	assert.Equal(t, daysToAdd(time.Monday, time.Monday), 7)
	assert.Equal(t, daysToAdd(time.Monday, time.Tuesday), 1)
	assert.Equal(t, daysToAdd(time.Tuesday, time.Monday), 6)
}